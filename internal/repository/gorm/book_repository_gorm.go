package gorm

import (
	"context"
	"errors"
	"fmt"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type bookRepositoryGORM struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) usecase.BookRepository {
	return &bookRepositoryGORM{db: db}
}

func (r *bookRepositoryGORM) Create(ctx context.Context, book *domain.Book) error {
	gormBook := BookFromDomain(book)
	result := r.db.WithContext(ctx).Create(gormBook)
	if result.Error != nil {
		return result.Error
	}
	book.ID = gormBook.ID
	return nil
}

func (r *bookRepositoryGORM) FindAll(ctx context.Context) ([]domain.Book, error) {
	var gormBooks []BookGORM
	if err := r.db.WithContext(ctx).Order("book_order asc, title asc").Find(&gormBooks).Error; err != nil {
		return nil, err
	}

	var domainBooks []domain.Book
	for _, b := range gormBooks {
		domainBooks = append(domainBooks, *b.ToDomain())
	}
	return domainBooks, nil
}

func (r *bookRepositoryGORM) FindByID(ctx context.Context, id uint) (*domain.Book, error) {
	var gormBook BookGORM
	result := r.db.WithContext(ctx).First(&gormBook, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormBook.ToDomain(), nil
}

func (r *bookRepositoryGORM) Update(ctx context.Context, book *domain.Book) error {
	gormBook := BookFromDomain(book)
	return r.db.WithContext(ctx).Save(gormBook).Error
}

func (r *bookRepositoryGORM) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&BookGORM{}, id).Error
}

func (r *bookRepositoryGORM) FindByOrder(ctx context.Context, order int) (*domain.Book, error) {
	var gormBook BookGORM
	result := r.db.WithContext(ctx).Where("book_order = ?", order).First(&gormBook)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormBook.ToDomain(), nil
}

func (r *bookRepositoryGORM) CreateBookWithOrderShift(ctx context.Context, book *domain.Book) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		newOrder := book.BookOrder
		// log.Printf("[DEBUG] CreateBookWithOrderShift: Attempting to insert at order %d for book '%s'", newOrder, book.Title)

		// 1. Find IDs of books to shift (ordered highest first)
		var idsToShift []uint
		if err := tx.Model(&BookGORM{}).
			Where("book_order >= ?", newOrder).
			Order("book_order DESC"). // Highest order first is CRITICAL here
			Pluck("id", &idsToShift).Error; err != nil {
			// log.Printf("[ERROR] Failed finding books to shift: %v", err)
			return err
		}
		// log.Printf("[DEBUG] Found %d books to shift (IDs: %v)", len(idsToShift), idsToShift)

		// 2. Shift them one by one, starting from the highest order
		if len(idsToShift) > 0 {
			// log.Printf("[DEBUG] Shifting books one by one...")
			for _, idToShift := range idsToShift {
				// log.Printf("[DEBUG] Incrementing book_order for ID: %d", idToShift)
				// Update book_order = book_order + 1 for the specific ID
				result := tx.Model(&BookGORM{}).Where("id = ?", idToShift).UpdateColumn("book_order", gorm.Expr("book_order + 1"))
				// We use UpdateColumn to avoid triggering hooks/UpdatedAt unnecessarily during the shift
				if result.Error != nil {
					// log.Printf("[ERROR] Failed shifting book ID %d: %v", idToShift, result.Error)
					return fmt.Errorf("failed shifting book ID %d: %w", idToShift, result.Error)
				}
				if result.RowsAffected == 0 {
					// log.Printf("[WARN] Shift update affected 0 rows for ID %d (maybe it was already updated/deleted?)", idToShift)
					// Decide if this should be a critical error or just a warning
				}
			}
			// log.Printf("[DEBUG] Finished shifting %d books individually.", len(idsToShift))
		} else {
			// log.Printf("[DEBUG] No books found at or after order %d, no shifting needed.", newOrder)
		}

		// 3. Create the new book
		gormBook := BookFromDomain(book)
		gormBook.BookOrder = newOrder
		// log.Printf("[DEBUG] Creating new book '%s' with final order %d", gormBook.Title, gormBook.BookOrder)
		if err := tx.Create(gormBook).Error; err != nil {
			// log.Printf("[ERROR] Failed during create after shift: %v", err)
			return err
		}

		book.ID = gormBook.ID
		// log.Printf("[DEBUG] Successfully created book ID %d at order %d", book.ID, newOrder)
		return nil // Commit transaction
	})
}

func (r *bookRepositoryGORM) UpdateBookWithOrderShift(ctx context.Context, book *domain.Book, newOrder int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var currentBookInDB BookGORM
		if err := tx.Model(&BookGORM{}).Select("book_order").First(&currentBookInDB, book.ID).Error; err != nil {
			// log.Printf("[ERROR] UpdateShift: Failed to get current book order for ID %d: %v", book.ID, err)
			return errors.New("failed to find book being updated")
		}
		oldOrder := currentBookInDB.BookOrder

		if oldOrder == newOrder {
			// log.Printf("[DEBUG] UpdateShift: Order not changed for book ID %d. Performing standard update.", book.ID)
			// If order hasn't changed, just do a normal save.
			gormBook := BookFromDomain(book)
			gormBook.BookOrder = newOrder // Ensure order is set
			if err := tx.Save(gormBook).Error; err != nil {
				// log.Printf("[ERROR] Failed saving target book ID %d (no order change): %v", book.ID, err)
				return err
			}
			return nil // Commit early
		}

		// log.Printf("[DEBUG] UpdateBookWithOrderShift: Moving book ID %d from order %d to %d", book.ID, oldOrder, newOrder)

		// Pindahkan buku yang ditarget ke 'order' sementara (cth: 0)
		// Ini untuk "membebaskan" slot 'oldOrder' agar bisa diisi oleh buku lain.
		// Kita menggunakan 0 karena validasi domain adalah >= 1.
		// log.Printf("[DEBUG] UpdateShift: Temporarily moving book ID %d to order 0", book.ID)
		if err := tx.Model(&BookGORM{}).Where("id = ?", book.ID).UpdateColumn("book_order", 0).Error; err != nil {
			// log.Printf("[ERROR] UpdateShift: Failed to move target book (ID %d) to temp order: %v", book.ID, err)
			return fmt.Errorf("failed setting temp order for target book: %w", err)
		}

		// 1. Shift books to make space or fill gap
		if newOrder < oldOrder {
			// Moving Up (e.g., 5 -> 2). Books [newOrder, oldOrder-1] need +1
			var idsToShift []uint
			if err := tx.Model(&BookGORM{}).
				Where("book_order >= ? AND book_order < ? AND id != ?", newOrder, oldOrder, book.ID).
				Order("book_order DESC"). // Process highest first
				Pluck("id", &idsToShift).Error; err != nil {
				// log.Printf("[ERROR] UpdateShift (Up): Failed finding books to shift: %v", err)
				return err
			}
			// log.Printf("[DEBUG] Moving up: Found %d books to increment (IDs: %v)", len(idsToShift), idsToShift)
			for _, idToShift := range idsToShift {
				// log.Printf("[DEBUG] Moving up: Incrementing order for ID %d", idToShift)
				if err := tx.Model(&BookGORM{}).Where("id = ?", idToShift).UpdateColumn("book_order", gorm.Expr("book_order + 1")).Error; err != nil {
					// log.Printf("[ERROR] Failed shifting book ID %d up: %v", idToShift, err)
					return fmt.Errorf("failed shifting book ID %d up: %w", idToShift, err)
				}
			}

		} else { // newOrder > oldOrder
			// Moving Down (e.g., 2 -> 5). Books [oldOrder+1, newOrder] need -1
			var idsToShift []uint
			if err := tx.Model(&BookGORM{}).
				Where("book_order > ? AND book_order <= ? AND id != ?", oldOrder, newOrder, book.ID).
				Order("book_order ASC"). // Process lowest first
				Pluck("id", &idsToShift).Error; err != nil {
				// log.Printf("[ERROR] UpdateShift (Down): Failed finding books to shift: %v", err)
				return err
			}
			// log.Printf("[DEBUG] Moving down: Found %d books to decrement (IDs: %v)", len(idsToShift), idsToShift)
			for _, idToShift := range idsToShift {
				// log.Printf("[DEBUG] Moving down: Decrementing order for ID %d", idToShift)
				if err := tx.Model(&BookGORM{}).Where("id = ?", idToShift).UpdateColumn("book_order", gorm.Expr("book_order - 1")).Error; err != nil {
					// log.Printf("[ERROR] Failed shifting book ID %d down: %v", idToShift, err)
					return fmt.Errorf("failed shifting book ID %d down: %w", idToShift, err)
				}
			}
		}

		// 2. Update the target book (Title, Desc, AND the new BookOrder)
		gormBook := BookFromDomain(book)
		gormBook.BookOrder = newOrder // Explicitly set the final order
		// log.Printf("[DEBUG] Updating target book ID %d with final order %d", gormBook.ID, gormBook.BookOrder)

		// Use Save to update all fields passed in 'book', plus the explicitly set BookOrder
		if err := tx.Save(gormBook).Error; err != nil {
			// log.Printf("[ERROR] Failed saving target book ID %d: %v", book.ID, err)
			return err
		}

		// log.Printf("[DEBUG] Successfully updated book ID %d to order %d", book.ID, newOrder)
		return nil // Commit transaction
	})
}

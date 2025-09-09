package dto

import "github.com/ipincamp/go-edsa-api/internal/domain"

// SubmitInteractionRequest adalah DTO untuk request body saat user menyelesaikan interaksi
type SubmitInteractionRequest struct {
	BookID        string `json:"book_id" binding:"required"`
	PageID        string `json:"page_id" binding:"required"`
	InteractionID string `json:"interaction_id" binding:"required"`
	Answer        string `json:"answer"` // Jawaban dari user, bisa berupa JSON string atau lainnya
}

// BookStatus merepresentasikan status buku untuk seorang user
type BookStatus string

const (
	StatusLocked     BookStatus = "locked"
	StatusAvailable  BookStatus = "available"
	StatusCompleted  BookStatus = "completed"
	StatusInProgress BookStatus = "in_progress"
)

// BookResponse adalah DTO untuk daftar buku
type BookResponse struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	CoverImageURL string     `json:"cover_image_url"`
	Level         int        `json:"level"`
	Status        BookStatus `json:"status"` // Status buku (locked, available, completed)
}

// PageResponse adalah DTO untuk halaman dalam sebuah buku
type PageResponse struct {
	ID         string `json:"id"`
	PageNumber int    `json:"page_number"`
	// Tambahkan field lain yang dibutuhkan frontend, misal content_data
}

// BookDetailResponse adalah DTO untuk detail sebuah buku
type BookDetailResponse struct {
	ID                  string         `json:"id"`
	Title               string         `json:"title"`
	UserProgressStatus  string         `json:"user_progress_status"` // 'not_started', 'in_progress', 'completed'
	LastCompletedPageID *string        `json:"last_completed_page_id"`
	Pages               []PageResponse `json:"pages"`
}

// Mapper untuk mengubah domain.Page menjadi DTO
func ToPageListResponse(pages []domain.Page) []PageResponse {
	var response []PageResponse
	for _, p := range pages {
		response = append(response, PageResponse{
			ID:         p.ID,
			PageNumber: p.PageNumber,
		})
	}
	return response
}

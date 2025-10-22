package admin

import (
	"context"
	"errors"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type adminService struct {
	subjectRepo     usecase.SubjectRepository
	classRepo       usecase.ClassRepository
	groupRepo       usecase.GroupRepository
	bookRepo        usecase.BookRepository
	pageRepo        usecase.PageRepository
	interactionRepo usecase.InteractionRepository
}

func NewAdminService(
	subjectRepo usecase.SubjectRepository,
	classRepo usecase.ClassRepository,
	groupRepo usecase.GroupRepository,
	bookRepo usecase.BookRepository,
	pageRepo usecase.PageRepository,
	interactionRepo usecase.InteractionRepository,
) usecase.AdminService {
	return &adminService{
		subjectRepo:     subjectRepo,
		classRepo:       classRepo,
		groupRepo:       groupRepo,
		bookRepo:        bookRepo,
		pageRepo:        pageRepo,
		interactionRepo: interactionRepo,
	}
}

// --- Mapper DTO ---
// (Helper untuk konversi domain ke response DTO)

func toSubjectResponse(s *domain.Subject) *domain.SubjectResponse {
	return &domain.SubjectResponse{
		ID:   s.ID,
		Name: s.Name,
	}
}

func toClassResponse(c *domain.Class) *domain.ClassResponse {
	resp := &domain.ClassResponse{
		ID:        c.ID,
		Name:      c.Name,
		SubjectID: c.SubjectID,
	}
	if c.Subject.ID != 0 {
		resp.Subject = *toSubjectResponse(&c.Subject)
	}
	return resp
}

func toGroupResponse(g *domain.Group) *domain.GroupResponse {
	resp := &domain.GroupResponse{
		ID:      g.ID,
		Name:    g.Name,
		ClassID: g.ClassID,
	}
	if g.Class.ID != 0 {
		resp.Class = *toClassResponse(&g.Class)
	}
	return resp
}

func toBookResponse(b *domain.Book) *domain.BookResponse {
	return &domain.BookResponse{
		ID:            b.ID,
		Title:         b.Title,
		Description:   b.Description,
		CoverImageURL: b.CoverImageURL,
		Theme:         b.Theme,
		BookOrder:     b.BookOrder,
	}
}

func toPageResponse(p *domain.Page) *domain.PageResponse {
	resp := &domain.PageResponse{
		ID:              p.ID,
		BookID:          p.BookID,
		PageNumber:      p.PageNumber,
		NarrativeText:   p.NarrativeText,
		InstructionText: p.InstructionText,
	}
	if p.Book.ID != 0 {
		resp.Book = *toBookResponse(&p.Book)
	}
	return resp
}

func toInteractionResponse(i *domain.Interaction) *domain.InteractionResponse {
	resp := &domain.InteractionResponse{
		ID:     i.ID,
		PageID: i.PageID,
		Type:   i.Type,
		Config: i.Config,
	}
	if i.Page.ID != 0 {
		resp.Page = *toPageResponse(&i.Page)
	}
	return resp
}

// --- Subject Methods ---

func (s *adminService) CreateSubject(ctx context.Context, req *domain.CreateSubjectRequest) (*domain.SubjectResponse, error) {
	subject := &domain.Subject{
		Name: req.Name,
	}
	if err := s.subjectRepo.Create(ctx, subject); err != nil {
		return nil, err
	}
	return toSubjectResponse(subject), nil
}

func (s *adminService) GetAllSubjects(ctx context.Context) ([]domain.SubjectResponse, error) {
	subjects, err := s.subjectRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var responses []domain.SubjectResponse
	for _, subject := range subjects {
		responses = append(responses, *toSubjectResponse(&subject))
	}
	return responses, nil
}

func (s *adminService) GetSubjectByID(ctx context.Context, id uint) (*domain.SubjectResponse, error) {
	subject, err := s.subjectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if subject == nil {
		return nil, errors.New("subject not found")
	}
	return toSubjectResponse(subject), nil
}

func (s *adminService) UpdateSubject(ctx context.Context, id uint, req *domain.UpdateSubjectRequest) (*domain.SubjectResponse, error) {
	subject, err := s.subjectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if subject == nil {
		return nil, errors.New("subject not found")
	}

	subject.Name = req.Name
	if err := s.subjectRepo.Update(ctx, subject); err != nil {
		return nil, err
	}
	return toSubjectResponse(subject), nil
}

func (s *adminService) DeleteSubject(ctx context.Context, id uint) error {
	// TODO: Cek apakah ada class yang masih terikat sebelum hapus?
	// Untuk saat ini, biarkan DB constraint (OnDelete:RESTRICT) yang menangani
	return s.subjectRepo.Delete(ctx, id)
}

// --- Class Methods ---

func (s *adminService) CreateClass(ctx context.Context, req *domain.CreateClassRequest) (*domain.ClassResponse, error) {
	// Validasi apakah SubjectID ada
	_, err := s.subjectRepo.FindByID(ctx, req.SubjectID)
	if err != nil {
		return nil, errors.New("failed to check subject")
	}
	if err == nil { // Jika error-nya nil, berarti subject tidak ditemukan (seharusnya)
		// FIXME: Logika FindByID perlu diperjelas. Asumsi: FindByID mengembalikan error jika tidak ada.
		// Mari kita asumsikan FindByID mengembalikan (nil, nil) jika not found
	}
	// Asumsi: Kita butuh subject ada
	subject, err := s.subjectRepo.FindByID(ctx, req.SubjectID)
	if err != nil {
		return nil, err // Error DB
	}
	if subject == nil {
		return nil, errors.New("subject not found")
	}

	class := &domain.Class{
		Name:      req.Name,
		SubjectID: req.SubjectID,
	}
	if err := s.classRepo.Create(ctx, class); err != nil {
		return nil, err
	}

	class.Subject = *subject // Attach subject yang sudah di-fetch
	return toClassResponse(class), nil
}

func (s *adminService) GetAllClasses(ctx context.Context) ([]domain.ClassResponse, error) {
	classes, err := s.classRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var responses []domain.ClassResponse
	for _, class := range classes {
		responses = append(responses, *toClassResponse(&class))
	}
	return responses, nil
}

func (s *adminService) GetClassByID(ctx context.Context, id uint) (*domain.ClassResponse, error) {
	class, err := s.classRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("class not found")
	}
	return toClassResponse(class), nil
}

func (s *adminService) UpdateClass(ctx context.Context, id uint, req *domain.UpdateClassRequest) (*domain.ClassResponse, error) {
	// Cek class
	class, err := s.classRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("class not found")
	}

	// Cek subject baru
	subject, err := s.subjectRepo.FindByID(ctx, req.SubjectID)
	if err != nil {
		return nil, err
	}
	if subject == nil {
		return nil, errors.New("subject not found")
	}

	class.Name = req.Name
	class.SubjectID = req.SubjectID

	if err := s.classRepo.Update(ctx, class); err != nil {
		return nil, err
	}

	class.Subject = *subject
	return toClassResponse(class), nil
}

func (s *adminService) DeleteClass(ctx context.Context, id uint) error {
	return s.classRepo.Delete(ctx, id)
}

// --- Group Methods ---

func (s *adminService) CreateGroup(ctx context.Context, req *domain.CreateGroupRequest) (*domain.GroupResponse, error) {
	// Cek ClassID
	class, err := s.classRepo.FindByID(ctx, req.ClassID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("class not found")
	}

	group := &domain.Group{
		Name:    req.Name,
		ClassID: req.ClassID,
	}
	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, err
	}

	group.Class = *class // Attach class
	return toGroupResponse(group), nil
}

func (s *adminService) GetAllGroups(ctx context.Context) ([]domain.GroupResponse, error) {
	groups, err := s.groupRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var responses []domain.GroupResponse
	for _, group := range groups {
		responses = append(responses, *toGroupResponse(&group))
	}
	return responses, nil
}

func (s *adminService) GetGroupByID(ctx context.Context, id uint) (*domain.GroupResponse, error) {
	group, err := s.groupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}
	return toGroupResponse(group), nil
}

func (s *adminService) UpdateGroup(ctx context.Context, id uint, req *domain.UpdateGroupRequest) (*domain.GroupResponse, error) {
	group, err := s.groupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, errors.New("group not found")
	}

	class, err := s.classRepo.FindByID(ctx, req.ClassID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("class not found")
	}

	group.Name = req.Name
	group.ClassID = req.ClassID

	if err := s.groupRepo.Update(ctx, group); err != nil {
		return nil, err
	}

	group.Class = *class
	return toGroupResponse(group), nil
}

func (s *adminService) DeleteGroup(ctx context.Context, id uint) error {
	return s.groupRepo.Delete(ctx, id)
}

// --- Book Methods ---

func (s *adminService) CreateBook(ctx context.Context, req *domain.CreateBookRequest) (*domain.BookResponse, error) {
	book := &domain.Book{
		Title:         req.Title,
		Description:   req.Description,
		CoverImageURL: req.CoverImageURL,
		Theme:         req.Theme,
		BookOrder:     req.BookOrder,
	}
	if err := s.bookRepo.CreateBookWithOrderShift(ctx, book); err != nil {
		return nil, err
	}
	return toBookResponse(book), nil
}

func (s *adminService) GetAllBooks(ctx context.Context) ([]domain.BookResponse, error) {
	books, err := s.bookRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var responses []domain.BookResponse
	for _, b := range books {
		responses = append(responses, *toBookResponse(&b))
	}
	return responses, nil
}

func (s *adminService) GetBookByID(ctx context.Context, id uint) (*domain.BookResponse, error) {
	book, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, errors.New("book not found")
	}
	return toBookResponse(book), nil
}

func (s *adminService) UpdateBook(ctx context.Context, id uint, req *domain.UpdateBookRequest) (*domain.BookResponse, error) {
	// 1. Ambil data buku yang ada
	book, err := s.bookRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, errors.New("book not found")
	}

	// 2. Cek field mana yang di-update
	// 'book' adalah object domain yang kita fetch, kita update di memory
	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Description != nil {
		book.Description = *req.Description
	}
	if req.CoverImageURL != nil {
		book.CoverImageURL = *req.CoverImageURL
	}
	if req.Theme != nil {
		book.Theme = *req.Theme
	}

	// 3. Cek apakah urutan berubah
	newOrder := book.BookOrder // Default ke urutan lama
	orderChanged := false
	if req.BookOrder != nil && *req.BookOrder != book.BookOrder {
		newOrder = *req.BookOrder
		orderChanged = true
	}

	// 4. Panggil repository yang sesuai
	if orderChanged {
		// Panggil logic 'shift' (dari jawaban saya sebelumnya)
		// 'book' sudah berisi Title/Desc/Theme yang baru
		if err := s.bookRepo.UpdateBookWithOrderShift(ctx, book, newOrder); err != nil {
			return nil, err
		}
	} else {
		// Panggil update biasa (hanya Title/Desc/Theme, tanpa 'shift')
		if err := s.bookRepo.Update(ctx, book); err != nil {
			return nil, err
		}
	}

	// 5. Kembalikan data yang sudah di-merge
	book.BookOrder = newOrder // Pastikan ordernya update
	return toBookResponse(book), nil
}

func (s *adminService) DeleteBook(ctx context.Context, id uint) error {
	return s.bookRepo.Delete(ctx, id)
}

// --- Page Methods ---

func (s *adminService) CreatePage(ctx context.Context, bookID uint, req *domain.CreatePageRequest) (*domain.PageResponse, error) {
	// Cek apakah BookID ada
	book, err := s.bookRepo.FindByID(ctx, bookID)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, errors.New("book not found")
	}

	page := &domain.Page{
		BookID:          bookID,
		PageNumber:      req.PageNumber,
		NarrativeText:   req.NarrativeText,
		InstructionText: req.InstructionText,
	}
	if err := s.pageRepo.Create(ctx, page); err != nil {
		return nil, err
	}

	page.Book = *book // Attach book
	return toPageResponse(page), nil
}

func (s *adminService) GetAllPagesForBook(ctx context.Context, bookID uint) ([]domain.PageResponse, error) {
	pages, err := s.pageRepo.FindAllByBookID(ctx, bookID)
	if err != nil {
		return nil, err
	}
	var responses []domain.PageResponse
	for _, p := range pages {
		responses = append(responses, *toPageResponse(&p))
	}
	return responses, nil
}

func (s *adminService) GetPageByID(ctx context.Context, id uint) (*domain.PageResponse, error) {
	page, err := s.pageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, errors.New("page not found")
	}
	return toPageResponse(page), nil
}

func (s *adminService) UpdatePage(ctx context.Context, id uint, req *domain.UpdatePageRequest) (*domain.PageResponse, error) {
	page, err := s.pageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, errors.New("page not found")
	}

	page.PageNumber = req.PageNumber
	page.NarrativeText = req.NarrativeText
	page.InstructionText = req.InstructionText

	if err := s.pageRepo.Update(ctx, page); err != nil {
		return nil, err
	}
	// page sudah terisi data `Book` dari FindByID
	return toPageResponse(page), nil
}

func (s *adminService) DeletePage(ctx context.Context, id uint) error {
	return s.pageRepo.Delete(ctx, id)
}

// --- Interaction Methods ---

func (s *adminService) CreateInteraction(ctx context.Context, pageID uint, req *domain.CreateInteractionRequest) (*domain.InteractionResponse, error) {
	// Cek apakah PageID ada
	page, err := s.pageRepo.FindByID(ctx, pageID)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, errors.New("page not found")
	}

	interaction := &domain.Interaction{
		PageID: pageID,
		Type:   req.Type,
		Config: req.Config,
	}
	if err := s.interactionRepo.Create(ctx, interaction); err != nil {
		return nil, err
	}

	interaction.Page = *page // Attach page
	return toInteractionResponse(interaction), nil
}

func (s *adminService) GetAllInteractionsForPage(ctx context.Context, pageID uint) ([]domain.InteractionResponse, error) {
	interactions, err := s.interactionRepo.FindAllByPageID(ctx, pageID)
	if err != nil {
		return nil, err
	}
	var responses []domain.InteractionResponse
	for _, i := range interactions {
		responses = append(responses, *toInteractionResponse(&i))
	}
	return responses, nil
}

func (s *adminService) GetInteractionByID(ctx context.Context, id uint) (*domain.InteractionResponse, error) {
	interaction, err := s.interactionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if interaction == nil {
		return nil, errors.New("interaction not found")
	}
	return toInteractionResponse(interaction), nil
}

func (s *adminService) UpdateInteraction(ctx context.Context, id uint, req *domain.UpdateInteractionRequest) (*domain.InteractionResponse, error) {
	interaction, err := s.interactionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if interaction == nil {
		return nil, errors.New("interaction not found")
	}

	interaction.Type = req.Type
	interaction.Config = req.Config

	if err := s.interactionRepo.Update(ctx, interaction); err != nil {
		return nil, err
	}
	// interaction sudah terisi data `Page.Book` dari FindByID
	return toInteractionResponse(interaction), nil
}

func (s *adminService) DeleteInteraction(ctx context.Context, id uint) error {
	return s.interactionRepo.Delete(ctx, id)
}

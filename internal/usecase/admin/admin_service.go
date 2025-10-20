package admin

import (
	"context"
	"errors"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type adminService struct {
	subjectRepo usecase.SubjectRepository
	classRepo   usecase.ClassRepository
	groupRepo   usecase.GroupRepository
}

func NewAdminService(
	subjectRepo usecase.SubjectRepository,
	classRepo usecase.ClassRepository,
	groupRepo usecase.GroupRepository,
) usecase.AdminService {
	return &adminService{
		subjectRepo: subjectRepo,
		classRepo:   classRepo,
		groupRepo:   groupRepo,
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

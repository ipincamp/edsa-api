package service

import (
	"context"
	"errors"
	"time"

	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/dto"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/repository"
	"github.com/ipincamp/go-edsa-api/pkg/cache"
	"gorm.io/gorm"
)

var (
	ErrGroupNotFound       = errors.New("course group not found")
	ErrRequestNotFound     = errors.New("join request not found")
	ErrAlreadyProcessed    = errors.New("request has already been processed")
	ErrNotAuthorized       = errors.New("not authorized to perform this action")
	ErrRoleStudentNotFound = errors.New("role 'student' not found")
	ErrAlreadyEnrolled     = errors.New("user is already enrolled in a class")
	ErrRequestExists       = errors.New("join request already exists for this class")
)

// CourseService adalah kontrak untuk service yang mengelola kelas dan siswa
type CourseService interface {
	ApplyToJoinGroup(ctx context.Context, userID, groupCode string) error
	HandleJoinRequest(ctx context.Context, teacherID, requestID string, approved bool) error
	GetMyClasses(ctx context.Context, teacherID string) ([]dto.CourseGroupResponse, error)
	GetStudentsByClass(ctx context.Context, groupID string) ([]dto.UserListResponse, error)
	GetTeachersByClass(ctx context.Context, groupID string) ([]dto.UserListResponse, error)
}

type courseService struct {
	db              *gorm.DB
	courseGroupRepo repository.CourseGroupRepository
	joinRequestRepo repository.JoinRequestRepository
	enrollmentRepo  repository.EnrollmentRepository
	userRepo        repository.UserRepository
	roleRepo        repository.RoleRepository
}

// NewCourseService membuat instance baru CourseService
func NewCourseService(
	db *gorm.DB,
	courseGroupRepo repository.CourseGroupRepository,
	joinRequestRepo repository.JoinRequestRepository,
	enrollmentRepo repository.EnrollmentRepository,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
) CourseService {
	return &courseService{
		db:              db,
		courseGroupRepo: courseGroupRepo,
		joinRequestRepo: joinRequestRepo,
		enrollmentRepo:  enrollmentRepo,
		userRepo:        userRepo,
		roleRepo:        roleRepo,
	}
}

// ApplyToJoinGroup adalah logika untuk guest mengajukan diri masuk kelas
func (s *courseService) ApplyToJoinGroup(ctx context.Context, userID string, groupCode string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	type result struct {
		err error
	}
	resultChan := make(chan result, 1)

	go func() {
		// 1. Cari kelas berdasarkan group code
		group, err := s.courseGroupRepo.FindByCode(ctx, groupCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				resultChan <- result{err: ErrGroupNotFound}
				return
			}
			resultChan <- result{err: err}
			return
		}

		// TODO: Tambahkan validasi di repository untuk mengecek apakah user sudah punya request atau sudah terdaftar.
		// Untuk saat ini, kita lakukan di sini.
		// Cek apakah user sudah mengirim request ke kelas ini atau sudah terdaftar.

		// 2. Buat permintaan baru
		newRequest := domain.JoinGroupRequest{
			ApplicantUserID: userID,
			TargetGroupID:   group.ID,
			Status:          "pending",
			RequestDate:     time.Now(),
		}

		if err := s.joinRequestRepo.Create(ctx, &newRequest); err != nil {
			// Handle jika ada constraint unique (user_id, group_id)
			resultChan <- result{err: err}
			return
		}

		resultChan <- result{err: nil}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case res := <-resultChan:
		return res.err
	}
}

// HandleJoinRequest adalah logika untuk guru menerima/menolak permintaan
func (s *courseService) HandleJoinRequest(ctx context.Context, teacherID string, requestID string, approved bool) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	type result struct {
		err error
	}
	resultChan := make(chan result, 1)

	go func() {
		// 1. Ambil data request
		request, err := s.joinRequestRepo.FindByID(ctx, requestID)
		if err != nil {
			resultChan <- result{err: ErrRequestNotFound}
			return
		}
		if request.Status != "pending" {
			resultChan <- result{err: ErrAlreadyProcessed}
			return
		}

		// TODO: 2. Verifikasi apakah teacherID berhak memproses request untuk kelas ini.
		// Ini memerlukan penambahan method di repository, misal: isTeacherInGroup(teacherID, groupID)
		// Untuk sekarang, kita asumsikan guru berhak.

		if !approved {
			// Jika ditolak, cukup update status dan selesai.
			err := s.joinRequestRepo.UpdateStatus(ctx, requestID, "rejected")
			resultChan <- result{err: err}
			return
		}

		// Jika diterima, jalankan logika dalam satu transaksi database
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// Buat instance repository dengan transaksi agar semua operasi dalam satu scope
			txJoinReqRepo := repository.NewJoinRequestRepository(tx)
			// txUserRepo := repository.NewUserRepository(tx)
			txEnrollmentRepo := repository.NewEnrollmentRepository(tx)

			// a. Update status permintaan menjadi 'approved'
			if err := txJoinReqRepo.UpdateStatus(ctx, requestID, "approved"); err != nil {
				return err
			}

			// b. Ambil role 'student' dari cache atau database
			studentRole, found := cache.GetRoleByName(constant.RoleStudent.String())
			if !found {
				// Fallback ke database jika tidak ada di cache
				dbRole, err := s.roleRepo.FindByName(ctx, constant.RoleStudent.String())
				if err != nil {
					return ErrRoleStudentNotFound
				}
				studentRole = *dbRole
			}

			// c. Ubah role user dari 'guest' menjadi 'student'
			userToUpdate := domain.User{ID: request.ApplicantUserID}
			if err := tx.Model(&userToUpdate).Update("role_id", studentRole.ID).Error; err != nil {
				return err
			}

			// d. Tambahkan user ke tabel enrollment
			enrollment := &domain.Enrollment{
				StudentID:         request.ApplicantUserID,
				CourseGroupID:     request.TargetGroupID,
				HomeroomTeacherID: teacherID, // Guru yang menyetujui menjadi wali kelas
				JoinDate:          time.Now(),
			}
			if err := txEnrollmentRepo.Create(ctx, enrollment); err != nil {
				return err
			}

			return nil // Commit transaksi
		})

		// Jika transaksi berhasil, update cache untuk user yang rolenya berubah
		if err == nil {
			updatedUser, findErr := s.userRepo.FindByID(ctx, request.ApplicantUserID)
			if findErr == nil {
				cache.AddUserToCache(*updatedUser)
			}
		}

		resultChan <- result{err: err}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case res := <-resultChan:
		return res.err
	}
}

// GetMyClasses mengambil daftar kelas yang diajar oleh seorang guru
func (s *courseService) GetMyClasses(ctx context.Context, teacherID string) ([]dto.CourseGroupResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	type result struct {
		resp []dto.CourseGroupResponse
		err  error
	}
	resultChan := make(chan result, 1)

	go func() {
		groups, err := s.courseGroupRepo.GetByTeacherID(ctx, teacherID)
		if err != nil {
			resultChan <- result{resp: nil, err: err}
			return
		}

		// Konversi dari domain.CourseGroup ke dto.CourseGroupResponse
		response := dto.ToCourseGroupListResponse(groups)
		resultChan <- result{resp: response, err: nil}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

// GetStudentsByClass mengambil daftar siswa dalam sebuah kelas
func (s *courseService) GetStudentsByClass(ctx context.Context, groupID string) ([]dto.UserListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	type result struct {
		resp []dto.UserListResponse
		err  error
	}
	resultChan := make(chan result, 1)

	go func() {
		students, err := s.userRepo.GetStudentsByGroupID(ctx, groupID)
		if err != nil {
			resultChan <- result{resp: nil, err: err}
			return
		}

		// Konversi dari domain.User ke dto.UserListResponse
		response := dto.ToUserListResponse(students)
		resultChan <- result{resp: response, err: nil}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

// GetTeachersByClass mengambil daftar guru yang mengajar di sebuah kelas
func (s *courseService) GetTeachersByClass(ctx context.Context, groupID string) ([]dto.UserListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	type result struct {
		resp []dto.UserListResponse
		err  error
	}
	resultChan := make(chan result, 1)

	go func() {
		teachers, err := s.userRepo.GetTeachersByGroupID(ctx, groupID)
		if err != nil {
			resultChan <- result{resp: nil, err: err}
			return
		}

		// Konversi dari domain.User ke dto.UserListResponse
		response := dto.ToUserListResponse(teachers)
		resultChan <- result{resp: response, err: nil}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

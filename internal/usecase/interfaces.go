package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
)

// --- Repositories ---

// RoleRepository mendefinisikan kontrak untuk akses data peran pengguna
type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*domain.Role, error)
}

// UserRepository mendefinisikan kontrak untuk persistensi data pengguna
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

// SubjectRepository mendefinisikan kontrak untuk data mata pelajaran
type SubjectRepository interface {
	Create(ctx context.Context, subject *domain.Subject) error
	FindAll(ctx context.Context) ([]domain.Subject, error)
	FindByID(ctx context.Context, id uint) (*domain.Subject, error)
	Update(ctx context.Context, subject *domain.Subject) error
	Delete(ctx context.Context, id uint) error
}

// ClassRepository mendefinisikan kontrak untuk data kelas
type ClassRepository interface {
	Create(ctx context.Context, class *domain.Class) error
	FindAll(ctx context.Context) ([]domain.Class, error)
	FindByID(ctx context.Context, id uint) (*domain.Class, error)
	Update(ctx context.Context, class *domain.Class) error
	Delete(ctx context.Context, id uint) error
}

// GroupRepository mendefinisikan kontrak untuk data grup
type GroupRepository interface {
	Create(ctx context.Context, group *domain.Group) error
	FindAll(ctx context.Context) ([]domain.Group, error)
	FindByID(ctx context.Context, id uint) (*domain.Group, error)
	Update(ctx context.Context, group *domain.Group) error
	Delete(ctx context.Context, id uint) error
}

// --- Services ---

// UserService mendefinisikan logika bisnis untuk pengguna
type UserService interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error)
}

// PasswordService mendefinisikan kontrak untuk hashing password
type PasswordService interface {
	Hash(password string) (string, error)
	Compare(password, hash string) (bool, error)
}

// TokenService mendefinisikan kontrak untuk pembuatan & validasi token
type TokenService interface {
	CreateToken(user *domain.User, duration time.Duration) (string, error)
	ValidateToken(tokenString string) (uuid.UUID, error)
}

// AdminService mendefinisikan logika bisnis untuk fitur admin
type AdminService interface {
	// Subject
	CreateSubject(ctx context.Context, req *domain.CreateSubjectRequest) (*domain.SubjectResponse, error)
	GetAllSubjects(ctx context.Context) ([]domain.SubjectResponse, error)
	GetSubjectByID(ctx context.Context, id uint) (*domain.SubjectResponse, error)
	UpdateSubject(ctx context.Context, id uint, req *domain.UpdateSubjectRequest) (*domain.SubjectResponse, error)
	DeleteSubject(ctx context.Context, id uint) error

	// Class
	CreateClass(ctx context.Context, req *domain.CreateClassRequest) (*domain.ClassResponse, error)
	GetAllClasses(ctx context.Context) ([]domain.ClassResponse, error)
	GetClassByID(ctx context.Context, id uint) (*domain.ClassResponse, error)
	UpdateClass(ctx context.Context, id uint, req *domain.UpdateClassRequest) (*domain.ClassResponse, error)
	DeleteClass(ctx context.Context, id uint) error

	// Group
	CreateGroup(ctx context.Context, req *domain.CreateGroupRequest) (*domain.GroupResponse, error)
	GetAllGroups(ctx context.Context) ([]domain.GroupResponse, error)
	GetGroupByID(ctx context.Context, id uint) (*domain.GroupResponse, error)
	UpdateGroup(ctx context.Context, id uint, req *domain.UpdateGroupRequest) (*domain.GroupResponse, error)
	DeleteGroup(ctx context.Context, id uint) error
}

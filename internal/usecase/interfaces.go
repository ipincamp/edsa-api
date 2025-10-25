package usecase

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
)

// --- Repositories ---

// RoleRepository mendefinisikan kontrak untuk akses data peran pengguna
type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*domain.Role, error)
	FindByID(ctx context.Context, id uint) (*domain.Role, error)
	FindAll(ctx context.Context) ([]domain.Role, error)
}

// UserRepository mendefinisikan kontrak untuk persistensi data pengguna
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
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
	FindGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Group, error)
}

// BookRepository mendefinisikan kontrak untuk data buku
type BookRepository interface {
	Create(ctx context.Context, book *domain.Book) error
	FindAll(ctx context.Context) ([]domain.Book, error)
	FindByID(ctx context.Context, id uint) (*domain.Book, error)
	Update(ctx context.Context, book *domain.Book) error
	Delete(ctx context.Context, id uint) error
	FindByOrder(ctx context.Context, order int) (*domain.Book, error)
	CreateBookWithOrderShift(ctx context.Context, book *domain.Book) error
	UpdateBookWithOrderShift(ctx context.Context, book *domain.Book, newOrder int) error
	FindPaginated(ctx context.Context, filters *domain.BookQuery) (*domain.PaginatedBooks, error)
}

// PageRepository mendefinisikan kontrak untuk data halaman
type PageRepository interface {
	Create(ctx context.Context, page *domain.Page) error
	FindAllByBookID(ctx context.Context, bookID uint) ([]domain.Page, error)
	FindByID(ctx context.Context, id uint) (*domain.Page, error)
	Update(ctx context.Context, page *domain.Page) error
	Delete(ctx context.Context, id uint) error
}

// InteractionRepository mendefinisikan kontrak untuk data interaksi
type InteractionRepository interface {
	Create(ctx context.Context, interaction *domain.Interaction) error
	FindAllByPageID(ctx context.Context, pageID uint) ([]domain.Interaction, error)
	FindByID(ctx context.Context, id uint) (*domain.Interaction, error)
	Update(ctx context.Context, interaction *domain.Interaction) error
	Delete(ctx context.Context, id uint) error
}

// UserBookProgressRepository mendefinisikan kontrak untuk data progres baca pengguna
type UserBookProgressRepository interface {
	FindOrCreate(ctx context.Context, progress *domain.UserBookProgress) error
	FindByUserAndBook(ctx context.Context, userID uuid.UUID, bookID uint) (*domain.UserBookProgress, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.UserBookProgress, error)
	Update(ctx context.Context, progress *domain.UserBookProgress) error
}

// GameRepository mendefinisikan kontrak untuk data game
type GameRepository interface {
	FindAll(ctx context.Context) ([]domain.Game, error)
	FindByID(ctx context.Context, id uint) (*domain.Game, error)
}

// UserGameScoreRepository mendefinisikan kontrak untuk data skor game pengguna
type UserGameScoreRepository interface {
	FindByUserAndGame(ctx context.Context, userID uuid.UUID, gameID uint) (*domain.UserGameScore, error)
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.UserGameScore, error)
	Upsert(ctx context.Context, score *domain.UserGameScore) error // Create atau Update
}

// ActivityLogRepository mendefinisikan kontrak untuk data log aktivitas pengguna
type ActivityLogRepository interface {
	Create(ctx context.Context, log *domain.ActivityLog) error
	CreateBatch(ctx context.Context, logs []domain.ActivityLog) error
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.ActivityLog, error)
	FindPaginatedByUserID(ctx context.Context, userID uuid.UUID, filters *domain.ActivityLogQuery) (*domain.PaginatedActivityLogs, error)
}

// MediaAssetRepository mendefinisikan kontrak untuk data aset media
type MediaAssetRepository interface {
	Create(ctx context.Context, asset *domain.MediaAsset) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.MediaAsset, error)
	SoftDelete(ctx context.Context, assetID uuid.UUID, deleterID *uuid.UUID) error
}

// GroupBookSettingRepository mendefinisikan kontrak untuk pengaturan buku grup
type GroupBookSettingRepository interface {
	Upsert(ctx context.Context, setting *domain.GroupBookSetting) error
	FindSettingsByGroupIDs(ctx context.Context, groupIDs []uint) ([]domain.GroupBookSetting, error)
}

// --- Services ---

// UserService mendefinisikan logika bisnis untuk pengguna
type UserService interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserResponse, error)
	RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.TokenResponse, error)
	Logout(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
	ChangePassword(ctx context.Context, userID uuid.UUID, req *domain.ChangePasswordRequest) error
	UpdateUserDetails(ctx context.Context, userID uuid.UUID, req *domain.UpdateDetailsRequest) (*domain.UserResponse, error)
	RequestAccountDeletion(ctx context.Context, userID uuid.UUID) error
	ConfirmAccountDeletion(ctx context.Context, userID uuid.UUID, req *domain.ConfirmDeletionRequest) error
}

// PasswordService mendefinisikan kontrak untuk hashing password
type PasswordService interface {
	Hash(password string) (string, error)
	Compare(password, hash string) (bool, error)
}

// TokenService mendefinisikan kontrak untuk pembuatan & validasi token
type TokenService interface {
	CreateToken(user *domain.User, sessionID uuid.UUID, duration time.Duration) (string, error)
	ValidateToken(tokenString string) (userID uuid.UUID, sessionID uuid.UUID, err error)
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

	// Book
	CreateBook(ctx context.Context, req *domain.CreateBookRequest) (*domain.BookResponse, error)
	GetAllBooks(ctx context.Context, filters *domain.BookQuery) (*domain.PaginatedDTO, error)
	GetBookByID(ctx context.Context, id uint) (*domain.BookResponse, error)
	UpdateBook(ctx context.Context, id uint, req *domain.UpdateBookRequest) (*domain.BookResponse, error)
	DeleteBook(ctx context.Context, id uint) error

	// Page
	CreatePage(ctx context.Context, bookID uint, req *domain.CreatePageRequest) (*domain.PageResponse, error)
	GetAllPagesForBook(ctx context.Context, bookID uint) ([]domain.PageResponse, error)
	GetPageByID(ctx context.Context, id uint) (*domain.PageResponse, error)
	UpdatePage(ctx context.Context, id uint, req *domain.UpdatePageRequest) (*domain.PageResponse, error)
	DeletePage(ctx context.Context, id uint) error

	// Interaction
	CreateInteraction(ctx context.Context, pageID uint, req *domain.CreateInteractionRequest) (*domain.InteractionResponse, error)
	GetAllInteractionsForPage(ctx context.Context, pageID uint) ([]domain.InteractionResponse, error)
	GetInteractionByID(ctx context.Context, id uint) (*domain.InteractionResponse, error)
	UpdateInteraction(ctx context.Context, id uint, req *domain.UpdateInteractionRequest) (*domain.InteractionResponse, error)
	DeleteInteraction(ctx context.Context, id uint) error
}

// Mendefinisikan logika bisnis untuk Modul "Read" & "Game"
type AppService interface {
	// Read Module
	GetBooksWithProgress(ctx context.Context, userID uuid.UUID) ([]domain.BookResponse, error)
	GetProgressToRestore(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, bookID uint) (*domain.RestoreProgressResponse, error)
	UpdatePageProgress(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, req *domain.UpdateProgressRequest) error
	CompleteBookProgress(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, req *domain.CompleteProgressRequest) error
	// Game Module
	GetAllGames(ctx context.Context, userID uuid.UUID) ([]domain.GameResponse, error)
	SubmitGameScore(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, gameID uint, req *domain.SubmitGameScoreRequest) error
}

// ActivityLoggerService mendefinisikan logika bisnis untuk pencatatan aktivitas pengguna
type ActivityLoggerService interface {
	Log(ctx context.Context, logData domain.ActivityLog)
	Shutdown(ctx context.Context) error
}

// Mendefinisikan logika bisnis untuk Dashboard Guru
type DashboardService interface {
	GetStudentActivity(ctx context.Context, studentID uuid.UUID, filters *domain.ActivityLogQuery) (*domain.PaginatedDTO, error)
	UnlockBookForGroup(ctx context.Context, groupID uint, bookID uint) error
	// TODO: Tambahkan method dashboard lainnya
}

// FileStorageService mendefinisikan kontrak untuk mengunggah file
type FileStorageService interface {
	Upload(file *multipart.FileHeader, fileID uuid.UUID) (filePath string, err error)
	Delete(filePath string) error
}

// MediaService mendefinisikan logika bisnis untuk manajemen media
type MediaService interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, ownerID, ownerType string, uploaderID *uuid.UUID) (*domain.MediaAssetResponse, error)
	DeleteFile(ctx context.Context, assetID uuid.UUID, deleterID *uuid.UUID) error
}

// SessionBlacklistService mendefinisikan kontrak untuk blacklist sesi
type SessionBlacklistService interface {
	BlacklistSession(ctx context.Context, sessionID uuid.UUID, duration time.Duration) error
	IsSessionBlacklisted(ctx context.Context, sessionID uuid.UUID) (bool, error)
}

// EmailService mendefinisikan kontrak untuk mengirim email
type EmailService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

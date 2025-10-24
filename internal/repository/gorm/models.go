package gorm

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleGORM adalah representasi tabel 'roles' di database
type RoleGORM struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"type:varchar(50);uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (RoleGORM) TableName() string {
	return "roles"
}

// UserGORM adalah representasi tabel 'users' di database
type UserGORM struct {
	ID        uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primarykey"`
	Name      string      `gorm:"type:varchar(255)"`
	Email     string      `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string      `gorm:"type:varchar(255);not null"`
	RoleID    uint        `gorm:"not null"`
	Role      RoleGORM    `gorm:"foreignKey:RoleID"`
	Groups    []GroupGORM `gorm:"many2many:user_groups;"` // Relasi many-to-many
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserGORM) TableName() string {
	return "users"
}

// SubjectGORM adalah representasi tabel 'subjects' di database
type SubjectGORM struct {
	ID        uint        `gorm:"primarykey"`
	Name      string      `gorm:"type:varchar(255);not null"`
	Classes   []ClassGORM `gorm:"foreignKey:SubjectID"` // Relasi one-to-many
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (SubjectGORM) TableName() string {
	return "subjects"
}

// ClassGORM adalah representasi tabel 'classes' di database
type ClassGORM struct {
	ID        uint        `gorm:"primarykey"`
	Name      string      `gorm:"type:varchar(255);not null"`
	SubjectID uint        `gorm:"not null"`
	Subject   SubjectGORM `gorm:"foreignKey:SubjectID"`
	Groups    []GroupGORM `gorm:"foreignKey:ClassID"` // Relasi one-to-many
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ClassGORM) TableName() string {
	return "classes"
}

// GroupGORM adalah representasi tabel 'groups' di database
type GroupGORM struct {
	ID        uint       `gorm:"primarykey"`
	Name      string     `gorm:"type:varchar(255);not null"`
	ClassID   uint       `gorm:"not null"`
	Class     ClassGORM  `gorm:"foreignKey:ClassID"`
	Users     []UserGORM `gorm:"many2many:user_groups;"` // Relasi many-to-many
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (GroupGORM) TableName() string {
	return "groups"
}

// UserGroupGORM adalah representasi tabel 'user_groups' di database
type UserGroupGORM struct {
	UserID  uuid.UUID `gorm:"primaryKey"`
	GroupID uint      `gorm:"primaryKey"`
}

func (UserGroupGORM) TableName() string {
	return "user_groups"
}

// BookGORM adalah representasi tabel 'books' di database
type BookGORM struct {
	ID            uint       `gorm:"primarykey"`
	Title         string     `gorm:"type:varchar(255);not null"`
	Description   string     `gorm:"type:text"`
	CoverImageURL string     `gorm:"type:varchar(255)"`
	Theme         string     `gorm:"type:varchar(100)"`
	BookOrder     int        `gorm:"default:0;uniqueIndex"`
	Pages         []PageGORM `gorm:"foreignKey:BookID"` // Relasi one-to-many
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (BookGORM) TableName() string {
	return "books"
}

// PageGORM adalah representasi tabel 'pages' di database
type PageGORM struct {
	ID              uint              `gorm:"primarykey"`
	BookID          uint              `gorm:"not null;index"`
	Book            BookGORM          `gorm:"foreignKey:BookID"`
	PageNumber      int               `gorm:"not null;index"`
	NarrativeText   string            `gorm:"type:text"`
	InstructionText string            `gorm:"type:text"`
	Interactions    []InteractionGORM `gorm:"foreignKey:PageID"` // Relasi one-to-many
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (PageGORM) TableName() string {
	return "pages"
}

// InteractionGORM adalah representasi tabel 'interactions' di database
type InteractionGORM struct {
	ID        uint            `gorm:"primarykey"`
	PageID    uint            `gorm:"not null;index"`
	Page      PageGORM        `gorm:"foreignKey:PageID"`
	Type      string          `gorm:"type:varchar(50);not null"`
	Config    json.RawMessage `gorm:"type:jsonb"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (InteractionGORM) TableName() string {
	return "interactions"
}

// UserBookProgressGORM adalah representasi tabel 'user_book_progress' di database
type UserBookProgressGORM struct {
	ID                   uint      `gorm:"primarykey"`
	UserID               uuid.UUID `gorm:"not null;uniqueIndex:idx_user_book"`
	User                 UserGORM  `gorm:"foreignKey:UserID"`
	BookID               uint      `gorm:"not null;uniqueIndex:idx_user_book"`
	Book                 BookGORM  `gorm:"foreignKey:BookID"`
	Status               string    `gorm:"type:varchar(50);default:'locked'"`
	HighestScore         int       `gorm:"default:0"`
	LastPageID           uint      `gorm:"default:0"`
	CurrentSessionPoints int       `gorm:"default:0"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt `gorm:"index"`
}

func (UserBookProgressGORM) TableName() string {
	return "user_book_progress"
}

// GameGORM adalah representasi tabel 'games' di database
type GameGORM struct {
	ID               uint   `gorm:"primarykey"`
	Name             string `gorm:"type:varchar(255);not null"`
	Type             string `gorm:"type:varchar(100)"`
	RelatedBookTheme string `gorm:"type:varchar(100)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (GameGORM) TableName() string {
	return "games"
}

// UserGameScoreGORM adalah representasi tabel 'user_game_scores' di database
type UserGameScoreGORM struct {
	ID           uint      `gorm:"primarykey"`
	UserID       uuid.UUID `gorm:"not null;uniqueIndex:idx_user_game"`
	User         UserGORM  `gorm:"foreignKey:UserID"`
	GameID       uint      `gorm:"not null;uniqueIndex:idx_user_game"`
	Game         GameGORM  `gorm:"foreignKey:GameID"`
	HighestScore int       `gorm:"default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (UserGameScoreGORM) TableName() string {
	return "user_game_scores"
}

// ActivityLogGORM adalah representasi tabel 'activity_logs' di database
type ActivityLogGORM struct {
	ID             uint      `gorm:"primarykey"`
	UserID         uuid.UUID `gorm:"not null;index"`
	User           UserGORM  `gorm:"foreignKey:UserID"`
	SessionID      uuid.UUID `gorm:"not null;index"`
	Action         string    `gorm:"type:varchar(100);not null;index"`
	TimestampStart time.Time `gorm:"not null"`
	DurationMs     *int
	Details        json.RawMessage `gorm:"type:jsonb"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (ActivityLogGORM) TableName() string {
	return "activity_logs"
}

// MediaAssetGORM adalah representasi tabel 'media_assets' di database
type MediaAssetGORM struct {
	ID               uuid.UUID  `gorm:"type:uuid;primarykey"`
	FileName         string     `gorm:"type:varchar(255);not null"`
	FilePath         string     `gorm:"type:varchar(255);not null"`
	PublicURL        string     `gorm:"type:varchar(255);not null"`
	MimeType         string     `gorm:"type:varchar(100)"`
	FileSize         int64      `gorm:"not null"`
	OwnerID          string     `gorm:"type:varchar(255);index"`
	OwnerType        string     `gorm:"type:varchar(100);index"`
	UploadedByUserID *uuid.UUID `gorm:"type:uuid;index"`
	UploadedByUser   UserGORM   `gorm:"foreignKey:UploadedByUserID"`
	DeletedByUserID  *uuid.UUID `gorm:"type:uuid;index"`
	DeletedByUser    UserGORM   `gorm:"foreignKey:DeletedByUserID"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (MediaAssetGORM) TableName() string {
	return "media_assets"
}

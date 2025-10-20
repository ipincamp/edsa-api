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
	BookOrder     int        `gorm:"default:0"`
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

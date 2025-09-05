package domain

import (
	"context"
	"encoding/json"
	"time"
)

type Role struct {
	ID          string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string          `gorm:"type:varchar(50);uniqueIndex;not null"`
	Permissions json.RawMessage `gorm:"type:jsonb"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RoleRepository interface {
	List(ctx context.Context) ([]Role, error)
	FindByID(ctx context.Context, id string) (Role, error)
	FindByName(ctx context.Context, name string) (Role, error)
	Update(ctx context.Context, role *Role) error
}

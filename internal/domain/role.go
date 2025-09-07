package domain

import (
	"encoding/json"
	"time"
)

// Role adalah entitas untuk data role user
type Role struct {
	ID          string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string          `gorm:"type:varchar(50);uniqueIndex;not null"`
	Permissions json.RawMessage `gorm:"type:jsonb"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GetPermissions parses the JSON permissions field and returns a map[string]bool
func (r *Role) GetPermissions() (map[string]bool, error) {
	if r.Permissions == nil {
		return make(map[string]bool), nil
	}
	var permissions map[string]bool
	err := json.Unmarshal(r.Permissions, &permissions)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

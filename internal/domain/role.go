package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Role struct {
	ID          string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string          `gorm:"type:varchar(50);uniqueIndex;not null"`
	Permissions json.RawMessage `gorm:"type:jsonb"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Permissions map[string]bool

func (p Permissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Permissions) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &p)
}

// GetPermissions parses the JSON permissions field and returns a map
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

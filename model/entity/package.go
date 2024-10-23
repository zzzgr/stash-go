package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Package struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	ActivityId  uint   `gorm:"type:int(11);not null;comment:'活动ID'" json:"activityId"`
	AccountName string `gorm:"type:text;not null;comment:'提取出的账号面名称'" json:"accountName"`

	Url     string  `json:"url"`
	Method  string  `json:"method"`
	Host    string  `json:"host"`
	Path    string  `json:"path"`
	Queries JSONMap `gorm:"type:text" json:"queries"`
	Headers JSONMap `gorm:"type:text" json:"headers"`
	Body    string  `json:"body"`

	IP string `json:"ip"`
}

// JSONMap is a custom type that implements the sql.Scanner and driver.Valuer interfaces
type JSONMap map[string]string

// Value converts the map to JSON for storage in the database
func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return "{}", nil
	}
	return json.Marshal(m)
}

// Scan converts JSON from the database back into a map
func (m *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON value: %v", value)
	}

	return json.Unmarshal(bytes, m)
}

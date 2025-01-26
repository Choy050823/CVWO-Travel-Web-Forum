package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Thread struct {
	ID             int         `json:"id"`
	Title          string      `json:"title"`
	Content        string      `json:"content"`
	AttachedImages StringArray `json:"attachedImages"` // New field
	UserID         int         `json:"userId"`
	CategoryID     int         `json:"categoryId"`
	CreatedAt      time.Time   `json:"createdAt"`
	Likes          int         `json:"likes"`
	LastActive     string      `json:"lastActive"`
}

// Add StringArray type for PostgreSQL array handling
type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil // Return empty array literal for PostgreSQL
	}
	// Convert to PostgreSQL array format: {value1, value2, value3}
	return "{" + strings.Join(a, ",") + "}", nil
}

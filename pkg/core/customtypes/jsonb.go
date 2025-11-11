package customTypes

import (
	"database/sql/driver"
	"encoding/json"
)

// JSONB is a custom type for PostgreSQL JSONB fields
type JSONB map[string]interface{}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return json.Marshal(map[string]interface{}{})
	}
	return json.Marshal(j)
}

// String returns the JSON string representation
func (j JSONB) String() string {
	bytes, err := json.Marshal(j)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

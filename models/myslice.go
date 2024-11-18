package models

import (
	"encoding/json"
	"fmt"
)

// self slice
type mySlice[T int | uint | float32 | string] []T

// go data -> db data
func (m mySlice[T]) Value() ([]byte, error) {
	return json.Marshal(m)
}

// db data -> go data
func (m *mySlice[T]) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}
	return json.Unmarshal(data, m)
}

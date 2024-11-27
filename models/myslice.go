package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tswcbyy1107/dns-service/utils"
)

// self slice
type mySlice[T uint | string] []T

// go data -> db data, []slice -> x,y,z string
func (ms mySlice[T]) Value() (driver.Value, error) {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ms)), ","), "[]"), nil
}

// db data -> go data, x,y,z string -> []slice
func (ms *mySlice[T]) Scan(value interface{}) error {
	if value == nil {
		*ms = []T{}
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}
	// construct slice
	bytes := []byte{'['}
	bytes = append(bytes, data...)
	bytes = append(bytes, ']')
	return json.Unmarshal(bytes, ms)
}

func (m mySlice[T]) Del(value T) []T {
	newS := mySlice[T]{}
	for _, sub := range m {
		if value == sub {
			continue
		}
		newS = append(newS, sub)
	}
	return newS
}

func (m mySlice[T]) Contains(value T) bool {
	return utils.Contains[T](m, value)
}

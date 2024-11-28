package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/tswcbyy1107/dns-service/utils"
)

// self slice
type MySlice[T uint | string] []T

func (ms MySlice[T]) Len() int           { return len(ms) }
func (ms MySlice[T]) Swap(i, j int)      { ms[i], ms[j] = ms[j], ms[i] }
func (ms MySlice[T]) Less(i, j int) bool { return ms[i] < ms[j] }

// go data -> db data, []slice -> x,y,z string
func (ms MySlice[T]) Value() (driver.Value, error) {
	sort.Sort(ms)
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ms)), ","), "[]"), nil
}

// db data -> go data, x,y,z string -> []slice
func (ms *MySlice[T]) Scan(value interface{}) error {
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

func (m MySlice[T]) Del(value T) []T {
	newS := MySlice[T]{}
	for _, sub := range m {
		if value == sub {
			continue
		}
		newS = append(newS, sub)
	}
	return newS
}

func (m MySlice[T]) Contains(value T) bool {
	return utils.Contains[T](m, value)
}

package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/wws2058/intranet-dns/utils"
)

// self slice
type MySlice[T uint | string] []T

func (ms MySlice[T]) Len() int           { return len(ms) }
func (ms MySlice[T]) Swap(i, j int)      { ms[i], ms[j] = ms[j], ms[i] }
func (ms MySlice[T]) Less(i, j int) bool { return ms[i] < ms[j] }

// go data -> db data, []slice -> x,y,z string
// func (ms MySlice[T]) Value() (driver.Value, error) {
// 	sort.Sort(ms)
// 	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ms)), ","), "[]"), nil
// }

// go data -> db data, []slice -> x,y,z string
func (ms MySlice[T]) Value() (driver.Value, error) {
	if len(ms) == 0 {
		return nil, nil
	}
	sort.Sort(ms)
	bytes, err := json.Marshal(ms)
	if err != nil {
		return nil, err
	}
	str := strings.Trim(string(bytes), "[]")
	return []byte(str), nil
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
	if string(data) == "null" {
		return nil
	}
	// construct slice
	bytes := []byte{'['}
	bytes = append(bytes, data...)
	bytes = append(bytes, ']')
	err := json.Unmarshal(bytes, ms)
	return err
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

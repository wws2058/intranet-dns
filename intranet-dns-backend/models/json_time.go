package models

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"

	"github.com/wws2058/intranet-dns/utils"
)

type JsonTime time.Time

// MarshalJSON() ([]byte, error)
func (j JsonTime) MarshalJSON() ([]byte, error) {
	if j.IsZero() {
		return []byte(`""`), nil
	}
	formatStr := fmt.Sprintf("\"%s\"", j.ToTime().Format(utils.DefaultTimeFormat))
	return []byte(formatStr), nil
}

// UnmarshalJSON([]byte) error
func (j *JsonTime) UnmarshalJSON(data []byte) error {
	timeStr, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation(utils.DefaultTimeFormat, timeStr, time.Local)
	if err != nil {
		return err
	}
	*j = JsonTime(t)
	return nil
}

func (j JsonTime) ToTime() time.Time {
	return time.Time(j)
}

func (j JsonTime) String() string {
	return time.Time(j).Format(utils.DefaultTimeFormat)
}

func (j JsonTime) IsZero() bool {
	return time.Time(j).IsZero()
}

// go data -> db data
func (j JsonTime) Value() (driver.Value, error) {
	if j.IsZero() {
		return nil, nil
	}
	return j.String(), nil
}

// db data -> go data
func (j *JsonTime) Scan(v interface{}) error {
	dbData, ok := v.(time.Time)
	if ok {
		*j = JsonTime(dbData)
		return nil
	}
	return fmt.Errorf("can't cover %v to time.Time", v)
}

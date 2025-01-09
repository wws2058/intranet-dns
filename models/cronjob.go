package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const (
	HttpType  = "http"
	FuncType  = "function"
	ShellType = "shell"
)

type Cronjob struct {
	BaseModel

	Name        string      `gorm:"type:varchar(32) not null; index:idk_name" json:"name,omitempty"`       // task name
	Spec        string      `gorm:"type:varchar(32) not null;" json:"spec,omitempty"`                      // spec & s m h D M W
	Creator     string      `gorm:"type:varchar(32) not null; index:idk_creator" json:"creator,omitempty"` // user name
	Description string      `gorm:"type:varchar(256) not null" json:"description,omitempty"`               // task desc
	Started     bool        `gorm:"type:tinyint(1) default true" json:"started"`                           // task switch
	LastSucceed bool        `gorm:"type:tinyint(1) default true" json:"last_succeed"`                      // task last status
	TaskType    string      `gorm:"type:varchar(32) not null" json:"task_type,omitempty"`                  // task type & shell or http or function
	TaskArgs    Args        `gorm:"type:varchar(1024) default null" json:"task_args,omitempty"`            // task args
	History     TaskHistory `gorm:"type:varchar(2048) default null" json:"history,omitempty"`              // task running history, 5
}

type Args struct {
	Url          string `json:"url,omitempty"`           // http url
	FunctionName string `json:"function_name,omitempty"` // function name
}

func (r Args) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *Args) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}
	return json.Unmarshal(data, r)
}

type TaskHistory []TaskRecord

type TaskRecord struct {
	UID     string   `json:"uid"`               // history task uid
	Succeed bool     `json:"succeed"`           // history task succeed
	CallAt  JsonTime `json:"call_at,omitempty"` // history start at
	Error   string   `json:"error,omitempty"`   // history task err msg
}

func (h TaskHistory) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *TaskHistory) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("value is not []byte, value: %v", value)
	}
	err := json.Unmarshal(data, h)
	return err
}

// store the last five records
func (h *TaskHistory) Add(history TaskRecord) {
	*h = append(*h, history)
	storeLenInt := 5
	if len(*h) > storeLenInt {
		tmp := *h
		tmp = tmp[len(*h)-storeLenInt:]
		*h = tmp
	}
}

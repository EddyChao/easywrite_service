package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringList struct {
	Values []string
}

func (sl StringList) MarshalJSON() ([]byte, error) {
	if sl.Values == nil {
		sl.Values = make([]string, 0)
	}
	return json.Marshal(sl.Values)
}

func (sl *StringList) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &sl.Values)
	if sl.Values == nil {
		sl.Values = make([]string, 0)
	}
	return err
}

func (sl *StringList) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &sl.Values)
	case string:
		return json.Unmarshal([]byte(v), &sl.Values)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

func (sl StringList) Value() (driver.Value, error) {
	return json.Marshal(sl.Values)
}

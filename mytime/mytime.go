package mytime

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(ct.UnixMilli(), 10)), nil
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	var millis int64
	err := json.Unmarshal(data, &millis)
	if err != nil {
		return err
	}
	ct.Time = time.UnixMilli(millis)
	return nil
}

func (t CustomTime) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *CustomTime) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		return nil
	}
	t.Time = value.(time.Time)
	return nil
}

func (t CustomTime) String() string {
	return t.Format("2006-01-02 15:04:05")
}

type CustomNullTime struct {
	gorm.DeletedAt
}

func (ct *CustomNullTime) MarshalJSON() ([]byte, error) {
	milli := ct.Time.UnixMilli()
	if !ct.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(milli, 10)), nil
}

func (ct *CustomNullTime) UnmarshalJSON(data []byte) error {
	var millis int64
	err := json.Unmarshal(data, &millis)
	if err != nil {
		return err
	}
	ct.Time = time.UnixMilli(millis)
	return nil
}

package husky

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Int64 int64

func (v Int64) Value() (driver.Value, error) {
	return int64(v), nil
}

func (v Int64) Int64() int64 {
	return int64(v)
}

func (id *Int64) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	switch v := src.(type) {
	case []byte:
		tv, err := strconv.ParseInt(string(v), 10, 64)
		if err != nil {
			return err
		}
		*id = Int64(tv)
	case string:
		tv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		*id = Int64(tv)
	case int64:
		*id = Int64(v)
	default:
		return errors.New("类型转换错误")
	}
	return nil
}

func (v Int64) MarshalJSON() ([]byte, error) {
	return []byte("\"" + strconv.FormatInt(int64(v), 10) + "\""), nil
}

func (v *Int64) UnmarshalJSON(src []byte) error {
	s, err := int64UnquoteIfQuoted(src)
	if err != nil {
		return err
	}
	d, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*v = Int64(d)
	return nil
}

func (v Int64) String() string {
	return strconv.FormatInt(int64(v), 10)
}

func (v Int64) MarshalBinary() (data []byte, err error) {
	return []byte(v.String()), nil

}

func (v *Int64) UnmarshalBinary(data []byte) error {
	return v.UnmarshalJSON(data)
}

func NewInt64(v int64) Int64 {
	return Int64(v)
}

func int64UnquoteIfQuoted(value interface{}) (string, error) {
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return "", fmt.Errorf("could not convert value '%+v' to byte array of type '%T'",
			value, value)
	}
	if len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	return string(bytes), nil
}

type Json struct {
	Raw json.RawMessage
}

func NewJson() Json {
	return Json{}
}

func (v Json) Value() (driver.Value, error) {
	if len(v.Raw) == 0 {
		return nil, nil
	}
	return v.Raw.MarshalJSON()
}

func (v *Json) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, &v.Raw)
	return err
}

func (v Json) MarshalJSON() ([]byte, error) {
	return v.Raw.MarshalJSON()
}

func (v *Json) UnmarshalJSON(data []byte) error {
	return v.Raw.UnmarshalJSON(data)
}

func (v *Json) Parse(data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	v.UnmarshalJSON(bytes)
	return nil
}

func (v *Json) To(data any) error {
	return json.Unmarshal(v.Raw, data)
}

func ToJson(data any) (Json, error) {
	var result Json
	if err := result.Parse(data); err != nil {
		return result, err
	}
	return result, nil
}

type Paging struct {
	Page int64   `form:"page" json:"page"`
	Rows int64   `form:"rows" json:"rows"`
	Sort *string `form:"sort" json:"sort"`
}

func (p *Paging) Offset() int {
	if p.Page < 1 || p.Page > 1000 {
		return 0
	}
	return int((p.Page - 1)) * p.Limit()
}

func (p *Paging) Limit() int {
	if p.Rows <= 0 || p.Rows > 1000 {
		return 30
	}
	return int(p.Rows)
}

type PagingData[T any] struct {
	Total int64 `json:"total"`
	Rows  []*T  `json:"rows"`
}

func NewPagingData[T any](total int64, rows []*T) PagingData[T] {
	return PagingData[T]{
		Total: total,
		Rows:  rows,
	}
}

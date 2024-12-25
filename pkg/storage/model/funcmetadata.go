/*
@author: panfengguo
@since: 2024/12/25
@desc: desc
*/
package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

type FuncMetaData struct {
	gorm.Model
	FunctionName string           `gorm:"primaryKey;column:function_name;unique;index:idx_function_name"`
	ImageUrl     string           `gorm:"column:image_url"`
	ResSpecList  ResourceSpecList `gorm:"column:resource_spec_list;type:json"`
}

func (r *ResourceSpecList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan ResourceSpecList, expected []byte, got %T", value)
	}
	return json.Unmarshal(b, r)
}

func (r ResourceSpecList) Value() (driver.Value, error) {
	return json.Marshal(r)
}

type ResourceSpecList []ResourceSpec

type ResourceSpec struct {
	Cpu        int32
	Memory     int32
	MinReplica int32
	MaxReplica int32
	Family     string
	Revision   int32
}

func (FuncMetaData) TableName() string {
	return "func_metadata"
}

/*
@author: panfengguo
@since: 2024/11/9
@desc: desc
*/
package storage

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
	Cpu               int32
	Memory            int32
	MinReplica        int32
	MaxReplica        int32
	EnableAutoScaling bool
	ServiceName       string
	Family            string
	Revision          int32
}

func (FuncMetaData) TableName() string {
	return "func_metadata"
}

type FuncTaskData struct {
	gorm.Model
	TaskArn       string `gorm:"primaryKey;column:task_arn;unique;index:idx_task_arn"`
	ServiceName   string `gorm:"column:service_name;index:idx_service_name"`
	FunctionName  string `gorm:"column:function_name;index:idx_function_name"`
	PrivateIpv4   string `gorm:"column:private_ipv4"`
	PublicIpv4    string `gorm:"column:public_ipv4"`
	Cpu           int32  `gorm:"column:cpu"`
	Memory        int32  `gorm:"column:memory"`
	Family        string `gorm:"column:family"`
	Revision      int32  `gorm:"column:revision"`
	LastStatus    string `gorm:"column:last_status"`
	DesiredStatus string `gorm:"column:desired_status"`
	LaunchType    string `gorm:"column:launch_type"`
}

func (FuncTaskData) TableName() string {
	return "func_task_data"
}

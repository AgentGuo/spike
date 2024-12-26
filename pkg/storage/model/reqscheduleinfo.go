/*
@author: panfengguo
@since: 2024/12/25
@desc: desc
*/
package model

import "gorm.io/gorm"

type ReqScheduleInfo struct {
	gorm.Model
	ReqId                uint64 `gorm:"column:req_id;primaryKey"`
	FunctionName         string `gorm:"column:function_name;index:idx_function_name"`
	PlacedAwsServiceName string `gorm:"column:placed_aws_service_name;index:idx_aws_service_name"`
	RequiredCpu          int32  `gorm:"column:required_cpu"`
	RequiredMemory       int32  `gorm:"column:required_memory"`
}

func (ReqScheduleInfo) TableName() string {
	return "req_scheduler_info"
}

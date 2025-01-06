/*
@author: panfengguo
@since: 2024/12/25
@desc: desc
*/
package model

import (
	"github.com/AgentGuo/spike/pkg/constants"
	"gorm.io/gorm"
)

type FuncInstance struct {
	gorm.Model
	AwsServiceName string                 `gorm:"primaryKey;column:aws_service_name;index:idx_service_name"`
	AwsTaskArn     string                 `gorm:"column:aws_task_arn;index:idx_task_arn"`
	FunctionName   string                 `gorm:"column:function_name;index:idx_function_name"`
	Ipv4           string                 `gorm:"column:ipv4"`
	Cpu            int32                  `gorm:"column:cpu"`
	Memory         int32                  `gorm:"column:memory"`
	AwsFamily      string                 `gorm:"column:aws_family"`
	AwsRevision    int32                  `gorm:"column:aws_revision"`
	LastStatus     string                 `gorm:"column:last_status"`
	DesiredStatus  string                 `gorm:"column:desired_status"`
	LaunchType     constants.InstanceType `gorm:"column:launch_type;type:varchar(20)"`
}

func (FuncInstance) TableName() string {
	return "func_instance"
}

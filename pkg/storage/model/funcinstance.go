/*
@author: panfengguo
@since: 2024/12/25
@desc: desc
*/
package model

import "gorm.io/gorm"

type FuncInstance struct {
	gorm.Model
	AwsServiceName string `gorm:"primaryKey;column:aws_service_name;index:idx_service_name"`
	AwsTaskArn     string `gorm:"column:aws_task_arn;unique;index:idx_task_arn"`
	FunctionName   string `gorm:"column:function_name;index:idx_function_name"`
	PrivateIpv4    string `gorm:"column:private_ipv4"`
	PublicIpv4     string `gorm:"column:public_ipv4"`
	Cpu            int32  `gorm:"column:cpu"`
	Memory         int32  `gorm:"column:memory"`
	AwsFamily      string `gorm:"column:aws_family"`
	AwsRevision    int32  `gorm:"column:aws_revision"`
	LastStatus     string `gorm:"column:last_status"`
	DesiredStatus  string `gorm:"column:desired_status"`
	LaunchType     int32  `gorm:"column:launch_type"`
}

func (FuncInstance) TableName() string {
	return "func_instance"
}

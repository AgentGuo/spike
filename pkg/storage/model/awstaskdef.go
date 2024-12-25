/*
@author: panfengguo
@since: 2024/12/25
@desc: desc
*/
package model

import "gorm.io/gorm"

type AwsTaskDef struct {
	gorm.Model
	TaskFamily   string `gorm:"column:task_family;index:idx_task_family_revision,unique"`
	TaskRevision int32  `gorm:"column:task_revision;index:idx_task_family_revision,unique"`
	FunctionName string `gorm:"column:function_name;index:idx_func_cpu_mem_img,unique"`
	Cpu          int32  `gorm:"column:cpu;index:idx_func_cpu_mem_img,unique"`
	Memory       int32  `gorm:"column:memory;index:idx_func_cpu_mem_img,unique"`
	ImageUrl     string `gorm:"column:image_url;index:idx_func_cpu_mem_img,unique"`
}

func (AwsTaskDef) TableName() string {
	return "aws_task_def"
}

/*
@author: panfengguo
@since: 2024/11/9
@desc: desc
*/
package storage

import (
	"github.com/AgentGuo/spike/cmd/server/config"
	"github.com/AgentGuo/spike/pkg/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
)

// sudo docker run --name spike-mysql -e MYSQL_ROOT_PASSWORD=faaspassword -p 3306:3306 -d mysql:8.0.31
type MysqlClient struct {
	db     *gorm.DB
	logger *logrus.Logger
}

var (
	mysqlInitOnce sync.Once
)

func NewMysqlClient() *MysqlClient {
	var initErr error
	mysqlInitOnce.Do(func() {
		initErr = initMysql(config.GetConfig().MysqlDsn)
	})
	if initErr != nil {
		logger.GetLogger().Fatal(initErr)
	}

	db, err := gorm.Open(mysql.Open(config.GetConfig().MysqlDsn), &gorm.Config{})
	if err != nil {
		logger.GetLogger().Fatal(err)
	}
	return &MysqlClient{
		db:     db,
		logger: logger.GetLogger(),
	}
}

func initMysql(dsn string) error {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移模式
	if err := db.AutoMigrate(&FuncMetaData{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&FuncTaskData{}); err != nil {
		return err
	}

	return nil
}

func (m *MysqlClient) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// func metadata

func (m *MysqlClient) CreateFuncMetaData(data *FuncMetaData) error {
	return m.db.Create(data).Error
}

func (m *MysqlClient) HasFuncMetaDataByFunctionName(functionName string) (bool, error) {
	var data []FuncMetaData
	err := m.db.Where(map[string]interface{}{"function_name": functionName}).Find(&data).Error
	if err != nil {
		return false, err
	} else if len(data) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (m *MysqlClient) GetFuncMetaDataByFunctionName(functionName string) (*FuncMetaData, error) {
	data := &FuncMetaData{}
	return data, m.db.Where(map[string]interface{}{"function_name": functionName}).First(data).Error
}

func (m *MysqlClient) GetFuncMetaDataByCondition(condition map[string]interface{}) ([]FuncMetaData, error) {
	var data []FuncMetaData
	return data, m.db.Where(condition).Find(&data).Error
}

func (m *MysqlClient) GetFuncMetaData(id uint, data *FuncMetaData) error {
	return m.db.First(data, id).Error
}

func (m *MysqlClient) UpdateFuncMetaData(data *FuncMetaData) error {
	return m.db.Save(data).Error
}

func (m *MysqlClient) DeleteFuncMetaDataByFunctionName(functionName string) error {
	return m.db.Unscoped().Where(map[string]interface{}{"function_name": functionName}).Delete(&FuncMetaData{}).Error
}

func (m *MysqlClient) DeleteFuncMetaData(id uint) error {
	return m.db.Unscoped().Delete(&FuncMetaData{}, id).Error
}

// func task data

func (m *MysqlClient) CreateFuncTaskData(data *FuncTaskData) error {
	return m.db.Create(data).Error
}

func (m *MysqlClient) GetFuncTaskDataByCondition(condition map[string]interface{}) ([]FuncTaskData, error) {
	var data []FuncTaskData
	err := m.db.Where(condition).Find(&data).Error
	return data, err
}

func (m *MysqlClient) GetFuncTaskDataByFunctionName(functionName string) ([]FuncTaskData, error) {
	return m.GetFuncTaskDataByCondition(map[string]interface{}{"function_name": functionName})
}

func (m *MysqlClient) GetFuncTaskDataByServiceName(serviceName string) ([]FuncTaskData, error) {
	return m.GetFuncTaskDataByCondition(map[string]interface{}{"service_name": serviceName})
}

func (m *MysqlClient) GetFuncTaskDataByTaskArn(taskArn string) ([]FuncTaskData, error) {
	return m.GetFuncTaskDataByCondition(map[string]interface{}{"task_arn": taskArn})
}

func (m *MysqlClient) UpdateFuncTaskData(data *FuncTaskData) error {
	return m.db.Save(data).Error
}

func (m *MysqlClient) UpdateFuncTaskDataBatch(data []FuncTaskData) error {
	tx := m.db.Begin()
	for _, instance := range data {
		if err := tx.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Save(&instance).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (m *MysqlClient) DeleteFuncTaskDataServiceName(serviceName string) error {
	return m.db.Unscoped().Where(map[string]interface{}{"service_name": serviceName}).Delete(&FuncTaskData{}).Error
}

func (m *MysqlClient) DeleteFuncTaskDataFunctionName(functionName string) error {
	return m.db.Unscoped().Where(map[string]interface{}{"function_name": functionName}).Delete(&FuncTaskData{}).Error
}

func (m *MysqlClient) DeleteFuncTaskDataByCondition(condition map[string]interface{}) error {
	return m.db.Unscoped().Where(condition).Delete(&FuncTaskData{}).Error
}

func (m *MysqlClient) DeleteFuncTaskData(id uint) error {
	return m.db.Unscoped().Delete(&FuncTaskData{}, id).Error
}

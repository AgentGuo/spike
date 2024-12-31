/*
@author: panfengguo
@since: 2024/11/4
@desc: desc
*/
package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

var (
	configPath     string
	configInstance *ServerConfig
	configOnce     sync.Once
)

type ServerConfig struct {
	ServerIp          string   `yaml:"server_ip"`
	ServerPort        int      `yaml:"server_port"`
	MysqlDsn          string   `yaml:"mysql_dsn"`
	LogLevel          string   `yaml:"log_level"`
	LogToFile         bool     `yaml:"log_to_file"`
	AwsCluster        string   `yaml:"aws_cluster"`
	AwsSubnets        []string `yaml:"aws_subnets"`
	EC2Provider       string   `yaml:"ec2_provider"`
	AwsSecurityGroups []string `yaml:"aws_security_groups"`
	TaskRole          string   `yaml:"task_role"`
	DispatchTimeout   int      `yaml:"dispatch_timeout"`
	AutoScaleStep     int      `yaml:"auto_scale_step"`
	AutoScaleWindow   int      `yaml:"auto_scale_window"`
}

func SetConfigPath(path string) {
	configPath = path
	GetConfig()
}

func GetConfig() *ServerConfig {
	configOnce.Do(func() {
		configInstance = &ServerConfig{
			ServerIp:          "127.0.0.1",
			ServerPort:        13306,
			MysqlDsn:          "root:spikepassword@tcp(127.0.0.1:3306)/spike?charset=utf8mb4&parseTime=True&loc=Local",
			LogLevel:          "debug",
			LogToFile:         false,
			AwsCluster:        "spike_cluster_mini",
			AwsSubnets:        []string{"subnet-01930cb57dbc12f7e", "subnet-0c77aae8c226d039c", "subnet-02bd39d1f8b337c22"},
			EC2Provider:       "Infra-ECS-Cluster-spikeclustermini-d985e674-EC2CapacityProvider-FufGynLGFE0q",
			AwsSecurityGroups: []string{"sg-02221dbcd555d5277"},
			TaskRole:          "PixelsFaaSRole",
			DispatchTimeout:   20,
			AutoScaleStep:     5,
			AutoScaleWindow:   60,
		}
		if configPath != "" {
			if fileContent, e := os.ReadFile(configPath); e == nil {
				if e := yaml.Unmarshal(fileContent, configInstance); e != nil {
					log.Fatal("Failed to unmarshal config file: ", e)
				}
			} else {
				log.Fatal("Failed to read config file: ", e)
			}
		}
	})
	return configInstance
}

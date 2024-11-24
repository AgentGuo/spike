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
	AwsFassCluster    string   `yaml:"aws_fass_cluster"`
	AwsSubnets        []string `yaml:"aws_subnets"`
	AwsSecurityGroups []string `yaml:"aws_security_groups"`
}

func SetConfigPath(path string) {
	configPath = path
	GetConfig()
}

func GetConfig() *ServerConfig {
	configOnce.Do(func() {
		configInstance = &ServerConfig{
			ServerIp:       "127.0.0.1",
			ServerPort:     13306,
			MysqlDsn:       "root:faaspassword@tcp(127.0.0.1:3306)/faas",
			LogLevel:       "info",
			AwsFassCluster: "fass_cluster",
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

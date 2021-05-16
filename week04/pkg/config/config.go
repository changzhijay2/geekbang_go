package config

import (
	"github.com/spf13/viper"
)

type MysqlConf struct {
	DBName   string
	Username string
	Password string
	Host     string
	Port     int
}

type GrpcConf struct {
	Host string
	Port int
}

var (
	MysqlConfig MysqlConf
	GrpcConfig  GrpcConf
)

func InitConf(confPath string) {
	if confPath == "" {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("../../config")
	} else {
		viper.SetConfigFile(confPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := subParse("mysql", &MysqlConfig); err != nil {
		panic(err)
	}
	if err := subParse("grpc", &GrpcConfig); err != nil {
		panic(err)
	}
}

func subParse(key string, value interface{}) error {
	sub := viper.Sub(key)
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}

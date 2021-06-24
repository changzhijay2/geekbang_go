package conf

import (
	"github.com/spf13/viper"
)

type TCPConf struct {
	Host            string
	Port            int
	MaxConn         int
	RecvBufferSize  int
	SendBufferSize int
}

var (
	TCPConfig TCPConf
)

func InitConf(confPath string) {
	if confPath == "" {
		viper.SetConfigName("conf")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")
	} else {
		viper.SetConfigFile(confPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := subParse("tcp", &TCPConfig); err != nil {
		panic(err)
	}
}

func subParse(key string, value interface{}) error {
	sub := viper.Sub(key)
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}
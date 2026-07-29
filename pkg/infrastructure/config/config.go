package config

import "github.com/spf13/viper"

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Mysql  MysqlConfig  `mapstructure:"mysql"`
	Upload FileConfig   `mapstructure:"file"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type MysqlConfig struct {
	DB           string `mapstructure:"db"`
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	MaxOpenConns int    `mapstructure:"maxOpenConns"`
	MaxIdleConns int    `mapstructure:"maxIdleConns"`
}

type FileConfig struct {
	Path string `mapstructure:"path"`
	Host string `mapstructure:"host"`
}

func NewConfig() *Config {
	viperObj := viper.New()
	viperObj.AddConfigPath("config")
	viperObj.SetConfigName("config")
	viperObj.SetConfigType("yaml")

	if err := viperObj.ReadInConfig(); err != nil {
		panic("配置文件读取失败:" + err.Error())
	}

	conf := &Config{}

	if err := viperObj.Unmarshal(&conf); err != nil {
		panic("配置文件序列化失败:" + err.Error())
	}

	return conf
}

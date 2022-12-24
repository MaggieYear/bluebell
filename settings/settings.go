package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 全局变量，用来保存程序所有的配置信息
var (
	AppSettings   = &AppConfig{}
	LogSettings   = &LogConfig{}
	MysqlSettings = &MySQLConfig{}
	RedisSettings = &RedisConfig{}
	AuthSettings  = &AuthConfig{}
)

type AppConfig struct {
	Name      string `mapstructure: "name"`
	Mode      string `mapstructure: "mode"`
	Version   string `mapstructure: "version"`
	Host      string `mapstructure: "host"`
	Port      int    `mapstructure: "port"`
	StartTime string `mapstructure: "starttime"`
	MachineID uint16 `mapstructure: "machineid"`
}

type LogConfig struct {
	Level      string `mapstructure: "level"`
	Filename   string `mapstructure: "filename"`
	MaxSize    int    `mapstructure: "maxsize"`
	MaxAge     int    `mapstructure: "maxage"`
	MaxBackups int    `mapstructure: "maxbackups"`
}

type MySQLConfig struct {
	Host         string `mapstructure: "host"`
	Port         int    `mapstructure: "port"`
	User         string `mapstructure: "user"`
	Password     string `mapstructure: "password"`
	DbName       string `mapstructure: "dbname"`
	MaxOpenConns int    `mapstructure: "maxopenconns"`
	MaxIdleConns int    `mapstructure: "maxidleconns"`
}

type RedisConfig struct {
	Host     string `mapstructure: "host"`
	Port     int    `mapstructure: "port"`
	DbName   int    `mapstructure: "db"`
	Password string `mapstructure: "password"`
	PoolSize int    `mapstructure: "poolsize"`
}
type AuthConfig struct {
	JWTExpire int `mapstructure: "jwtexpire"`
}

func Init() (err error) {
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")
	//viper.SetConfigType("json")

	viper.SetConfigFile("./conf/config.yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("viper.ReadInConfig() failed, file not found, err: %v\n", err)
		} else {
			fmt.Printf("viper.ReadInConfig() failed, err: %v\n", err)
		}
	}

	// 把读取到的配置信息反序列化到变量中
	unmarchalSettingInfo()

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		// 配置文件修改后，需要重新反序列化
		unmarchalSettingInfo()
	})
	return
}

func unmarchalSettingInfo() {
	err := viper.Unmarshal(&AppSettings)
	printError(err)
	err = viper.UnmarshalKey("log", &LogSettings)
	printError(err)
	err = viper.UnmarshalKey("mysql", &MysqlSettings)
	printError(err)
	err = viper.UnmarshalKey("redis", &RedisSettings)
	printError(err)
	err = viper.UnmarshalKey("auth", &AuthSettings)
	printError(err)
}

func printError(err error) {
	if err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}
}

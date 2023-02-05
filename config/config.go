package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config contains all the configurations for the application
type TConf struct {
	LogLevel string `yaml:"loglevel"`
	Redis    struct {
		Addr string `yaml:"addr"`
	}
	Mysql struct {
		Addr         string `yaml:"addr"`
		User         string `yaml:"user"`
		PassWord     string `yaml:"password"`
		DataBase     string `yaml:"database"`
		MaxIdleConns int    `yaml:"maxidleconns"`
		MaxOpenConns int    `yaml:"maxopenconns"`
	} `yaml:"mysql"`
	Port string
	Nsq  struct {
		Addr        string `yaml:"addr"`
		MaxAttempts uint16 `yaml:"maxAttempts"`
	}
	Task struct {
		StationTaskCron   string `yaml:"stationTaskCron"`
		EmptyOffTaskCron  string `yaml:"emptyOffTaskCron"`
		InvertBoxTaskCron string `yaml:"invertBoxTaskCron"`
	}
	Elastic struct {
		Addr string
		User string
		Pass string
	}
}

var Conf TConf

// LoadConfig loads the configuration from file and environment variables
func LoadConfig(configFile string) (err error) {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("app")

	if err = viper.ReadInConfig(); err != nil {
		return err
	}

	if err = viper.Unmarshal(&Conf); err != nil {
		return err
	}

	//监听配置文件变化，默认每5s监听一次
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件发生变化：", e.Name)
		if err = viper.Unmarshal(&Conf); err != nil {
			fmt.Printf("配置文件变更失败：%s", err)
		}
	})

	//tick := time.Tick(time.Duration(10) * time.Second)
	//go func() {
	//	for {
	//		select {
	//		case <-tick:
	//			viper.WatchConfig()
	//			viper.OnConfigChange(func(e fsnotify.Event) {
	//				fmt.Println("配置文件发生变化：", e.Name)
	//				if err = viper.Unmarshal(&config); err != nil {
	//					fmt.Printf("配置文件变更失败：%s", err)
	//				}
	//			})
	//		}
	//	}
	//}()

	return err
}

// GetConfig 获取配置
func GetConfig() TConf {
	return Conf
}

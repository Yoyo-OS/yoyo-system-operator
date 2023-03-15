package settings

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: YSO is an utility which allows you to perform maintenance
	tasks on your Yoyo OS installation.
*/

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Updates UpdatesConfig `json:"updates"`
}

type UpdatesConfig struct {
	Schedule string `json:"schedule"`
}

var Cnf *Config

func init() {
	viper.AddConfigPath("/etc/yso/")
	viper.AddConfigPath("config/")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()

	if err != nil {
		panic("Config error!")
	}

	err = viper.Unmarshal(&Cnf)

	if err != nil {
		panic("Config error!\n" + err.Error())
	}
}

func GetConfig() *Config {
	return Cnf
}

func GetConfigKeys() []string {
	return viper.AllKeys()
}

func GetConfigValue(key string) interface{} {
	return viper.Get(key)
}

func SetConfigValue(key string, value interface{}) {
	if key == "updates.schedule" {
		if value != "daily" && value != "weekly" && value != "monthly" {
			fmt.Println("Invalid value for updates.schedule!")
			return
		}
	} else if key == "updates.smart" {
		if value != true && value != false && value != "true" && value != "false" {
			fmt.Println("Invalid value for updates.smart!")
			return
		}
	}

	if value == "true" {
		value = true
	} else if value == "false" {
		value = false
	}

	viper.Set(key, value)
}

func SaveConfig() error {
	return viper.WriteConfig()
}

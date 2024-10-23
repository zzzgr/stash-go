package config

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"os"
)

var Conf *Config

func SetUp() {

	err := ifConfig("config.yaml")
	if err != nil {
		log.Fatalf(fmt.Sprintf("初始化配置出错: %s", err.Error()))
		return
	}

	// 加载配置
	Conf = loadConfig()
}

func ifConfig(src string) error {
	file, err := os.Stat(src)
	if err != nil {
		log.Info("初始化配置文件")
		// config.yaml文件
		_, err = os.Create(src)
		if err != nil {
			return errors.New(fmt.Sprintf("初始化配置文件失败: %s", err))
		}

		// 默认配置
		c := &Config{}
		appYaml, err := yaml.Marshal(c)
		if err != nil {
			return err
		}

		err = os.WriteFile(src, appYaml, 0777)
		if err != nil {
			return errors.New(fmt.Sprintf("初始化配置文件失败: %s", err))
		}

		return nil
	}

	if file.IsDir() {
		return errors.New(fmt.Sprintf("%s 不是一个文件, 请检查!", src))
	} else {
		return nil
	}
}

func loadConfig() *Config {
	// 配置文件名称
	viper.SetConfigName("config")
	// 配置文件扩展名
	viper.SetConfigType("yaml")
	// 配置文件所在路径
	viper.AddConfigPath("./")
	// 查找并读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		// 处理读取配置文件的错误
		log.Panic(fmt.Sprintf("配置文件格式错误: %s", err.Error()))
	}

	// 配置信息绑定到结构体变量
	c := new(Config)
	err = viper.Unmarshal(&c)
	if err != nil {
		fmt.Printf("加载化配置文件失败: %v \n", err)
	}

	return c

}

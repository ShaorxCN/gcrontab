package config

import (
	"encoding/json"
	"flag"
	"os"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Configer 配置文件接口，需要实现初始化和重启功能。

type Configer interface {
	Init() error
	Restart() error
	Stop()
}

// Management 是配置管理。
type Management struct {
	Configs map[string]Configer
}

// DefaultManagement 管理 config 默认实例。
var DefaultManagement *Management

func init() {
	DefaultManagement = &Management{
		Configs: make(map[string]Configer),
	}
}

// Register 将需要的配置信息注册到配置管理。
func (cm *Management) Register(name string, c Configer) {
	cm.Configs[name] = c
}

// Init 初始化。
func (cm *Management) Init() error {
	// 优先读取环境变量配置，读取失败或者配置错误则读取文件配置
	ok := cm.getConfigFromENV()
	if !ok {
		if err := cm.loadConfigFile(); err != nil {
			logrus.WithField("error", err).Error("load config file failed")
			return err
		}
	}
	return nil
}

// Stop 停止服务
func (cm *Management) Stop() {
	for _, config := range cm.Configs {
		config.Stop()
	}
}

// loadConfigFile 加载本地配置文件。
func (cm *Management) loadConfigFile() error {
	filePath := flag.String("cf", "./config.json", "config path")
	flag.Parse()
	data, err := os.ReadFile(*filePath)
	if err != nil {
		return err
	}
	// 将 json 解码到注册的 config 中。
	configs := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &configs); err != nil {
		return err
	}
	for key, raw := range configs {
		if config, ok := cm.Configs[key]; ok {
			if err := json.Unmarshal(raw, config); err != nil {
				return err
			}
		}
	}
	return nil
}

// getConfigFromENV 从 env 获取配置信息。
func (cm *Management) getConfigFromENV() bool {
	envCount := 0
	// 通过反射找到字段对应的环境变量值
	for _, config := range cm.Configs {
		sType := reflect.TypeOf(config)
		sValue := reflect.ValueOf(config)
		length := sValue.Elem().NumField()
		for i := 0; i < length; i++ {
			envName := sType.Elem().Field(i).Tag.Get("env")
			envValue := os.Getenv(envName)
			if envValue != "" {
				envCount++
				switch sValue.Elem().Field(i).Kind() {
				case reflect.String:
					sValue.Elem().Field(i).SetString(envValue)
				case reflect.Int64, reflect.Int:
					intVal, err := strconv.ParseInt(envValue, 10, 64)
					if err != nil {
						logrus.Error(err)
					} else {
						sValue.Elem().Field(i).SetInt(intVal)
					}

				}
			}
		}
	}
	return envCount != 0
}

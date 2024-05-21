package config

import (
	"encoding/json"
	"flag"
	"os"
	"reflect"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Configer 配置文件接口，需要实现初始化重启以及终止功能。

type Configer interface {
	Init() error
	Restart() error
	Stop()
}

// Management 配置管理。
type Management struct {
	Configs map[string]Configer
}

// DefaultManagement  默认实例。
var DefaultManagement *Management

func init() {
	DefaultManagement = &Management{
		Configs: make(map[string]Configer),
	}
}

// Register 将需要的配置信息注册到配置管理。
func (cm *Management) Register(name string, c Configer) {
	if _, ok := cm.Configs[name]; ok {
		logrus.Infof("config [%s] exists,overwrite...", name)
	}
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

// loadConfigFile 加载本地配置文件 默认当前config.json
func (cm *Management) loadConfigFile() error {
	filePath := flag.String("cf", "./config.json", "config path")
	flag.Parse()
	logrus.Infoln("read local config  file", *filePath)
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
	var envCount int
	// 通过反射找到字段对应的环境变量值
	for name, config := range cm.Configs {
		envCount = 0
		cType := reflect.TypeOf(config)
		cValue := reflect.ValueOf(config)
		length := cValue.Elem().NumField()
		for i := 0; i < length; i++ {
			envName := cType.Elem().Field(i).Tag.Get("env")
			envValue := os.Getenv(envName)
			if envValue != "" {
				envCount++
				switch cValue.Elem().Field(i).Kind() {
				case reflect.String:
					cValue.Elem().Field(i).SetString(envValue)
				case reflect.Int, reflect.Int64:
					intVal, err := strconv.ParseInt(envValue, 10, 64)
					if err != nil {
						logrus.Error(err)
					} else {
						cValue.Elem().Field(i).SetInt(intVal)
					}
				}
			}
		}
		if envCount != length {
			logrus.WithField("configName", name).Info("init with default value")
		}
	}
	return envCount != 0
}

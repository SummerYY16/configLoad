package g

import (
	yaml "gopkg.in/yaml.v2"
	"sync/atomic"

	log "configLoad/toolkits/logger"
	"configLoad/toolkits/reconf"
)

// AppConfig 配置字段
type AppConfig struct {
	Apiserver struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
}

var appConfig AppConfig

// AppConfigMgr reload()协程写 和 for循环的读，都是对Appconfig对象，因此有读写冲突
type AppConfigMgr struct {
	Config atomic.Value
}

// AppConfigMgrHandler  初始化结构体
var AppConfigMgrHandler = &AppConfigMgr{}

// Callback 回调函数
func (a *AppConfigMgr) Callback(conf *reconf.Config) {
	err := yaml.Unmarshal(conf.Data, &appConfig)
	if err != nil {
		log.Error("解析yaml配置文件出错:", err)
	}

	AppConfigMgrHandler.Config.Store(&appConfig)
	// WebAPIRestart()
}

// InitConfig 初始化配置文件
func InitConfig(file string) {
	// [1] 打开配置文件
	conf, err := reconf.NewConfig(file)
	// conf, err := reconf.NewConfig("setting.yml")
	if err != nil {
		log.Error("read config file err: %v\n", err)
		return
	}

	// 添加观察者
	conf.AddObserver(AppConfigMgrHandler)

	// [2]第一次读取配置文件
	err = yaml.Unmarshal(conf.Data, &appConfig)
	if err != nil {
		log.Error("解析yaml配置文件出错:", err)
	}

	// [3] 把读取到的配置文件数据存储到atomic.Value
	AppConfigMgrHandler.Config.Store(&appConfig)
	log.Info("配置文件加载完成.")
}

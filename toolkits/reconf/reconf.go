package reconf

// 实现一个解析配置文件的包
import (
	"time"
	"os"
	"sync"
	"io/ioutil"

	log "configLoad/toolkits/logger"
)

// Config 配置文件字段
type Config struct {
	filename string
	Data []byte
	lastModifyTime int64
	rwLock sync.RWMutex
	notifyList []Notifyer
}

// NewConfig 配置文件解析
func NewConfig(file string) (conf *Config, err error) {
	curModifyTime, err := getLastModifyTime(file)
	if err != nil {
		log.Error("获取文件修改时间报错:%s\n", err)
	}
	conf = &Config{
		filename: file,
		lastModifyTime: curModifyTime,
	}

	yamlFile, err := ioutil.ReadFile(file)
    if err != nil {
        log.Error("yamlFile.Get err #%v ", err)
    }
	// 将解析配置文件后的数据更新到结构体的map中，写锁
	conf.rwLock.Lock()
	conf.Data = yamlFile
	conf.rwLock.Unlock()
	if err != nil {
        log.Error("Unmarshal: %v", err)
    }

	// 启一个后台线程去检测配置文件是否更改
	go conf.reload()
	return
}

// AddObserver 添加观察者
func (c *Config) AddObserver(n Notifyer) {
	c.notifyList = append(c.notifyList, n)
}

func (c *Config) reload(){
	// 定时器
	ticker := time.NewTicker(time.Second * 10) 
	for range ticker.C {
		func () {
			curModifyTime, err := getLastModifyTime(c.filename)
			if err != nil {
				log.Error("获取文件修改时间报错:%s\n", err)
			}
			if curModifyTime > c.lastModifyTime {
				yamlFile, err := ioutil.ReadFile(c.filename)
    			if err != nil {
    			    log.Error("yamlFile.Get err #%v ", err)
    			}
				// 将解析配置文件后的数据更新到结构体的map中，写锁
				c.rwLock.Lock()
				c.Data = yamlFile
				c.rwLock.Unlock()
				if err != nil {
    			    log.Error("Unmarshal: %v", err)
    			}

				c.lastModifyTime = curModifyTime

				// 配置更新通知所有观察者
				for _, n := range c.notifyList {
					n.Callback(c)
				}
			}
		}()
	}
}

func getLastModifyTime(file string) (lastModifyTime int64, err error) {
	f, err := os.Open(file)
	if err != nil {
		log.Error("打开文件报错:%s\n", err)
		return 
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		log.Error("获取配置文件状态错误:%s\n", err)
		return
	}
	// 或取当前文件修改时间
	curModifyTime := fileInfo.ModTime().Unix()
	return curModifyTime, nil
}

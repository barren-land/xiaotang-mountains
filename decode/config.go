package decode

import (
	"io/ioutil"
	"os"

	logs "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type ConfigInfo struct {
	Whitelist []string `yaml:"whitelist"`
	Port      int      `yaml:"port"`
}

func New() *ConfigInfo {
	return new(ConfigInfo)
}

func (configInfo *ConfigInfo) Parse(filePath string) {
	if filePath == "" {
		logs.WithFields(logs.Fields{
			"module": "解析配置文件",
		}).Fatalln("配置文件路径不能为空!")
		return
	}
	fd, err := os.Open(filePath)
	if err != nil {
		logs.WithFields(logs.Fields{
			"module": "解析配置文件",
		}).Fatalln("获取文件描述符异常,err:" + err.Error())
	}
	content, err := ioutil.ReadAll(fd)
	if err != nil {
		logs.WithFields(logs.Fields{
			"module": "解析配置文件",
		}).Fatalln("获取文件内容异常,err:" + err.Error())
	}
	if err := yaml.Unmarshal(content, configInfo); err != nil {
		logs.WithFields(logs.Fields{
			"module": "解析配置文件",
		}).Fatalln("yaml格式化文件异常,err:" + err.Error())
	}
}

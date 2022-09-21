package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/sluggard/poc/util"
	"github.com/spf13/viper"
)

type Config struct {
	Server      Server
	DevicesInfo DevicesInfo
	RunConfig   RunConfig
}

type Server struct {
	Host        string
	Port        int
	ContextPath string
}

type RunConfig struct {
	Remote bool
}

type DevicesInfo struct {
	Username string
	Password string
	Hosts    []string
}

const (
	DefaultConfigPath string = "conf/config.json"
)

var config Config = Config{
	Server: Server{
		Host:        "0.0.0.0",
		Port:        5678,
		ContextPath: "",
	},
	DevicesInfo: DevicesInfo{
		Username: "root",
		Password: "password",
		Hosts:    []string{"192.168.1.2", "192.168.1.3", "192.168.1.4", "192.168.1.5", "192.168.1.6"},
	},
	RunConfig: RunConfig{
		Remote: false,
	},
}

func GetConfig() *Config {
	return &config
}

func New(config string) Config {
	// return Config{"127.0.0.1", 5678}
	return LoadConfing(DefaultConfigPath)
}

func LoadConfing(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v", err)
	} else {
		//读取的数据为json格式，需要进行解码
		err = json.Unmarshal(data, &config)
		if err != nil {
			panic(err)
		}
	}
	return config
}

func MakeDemoFile(host string) {
	home, err := util.GetHome()
	log.Debug(home)
	if err != nil {
		return
	}
	fileName := "/root/config1.properties"
	dir, name := filepath.Split(fmt.Sprintf("%s%s.poc%s.%s%s", home, string(filepath.Separator), string(filepath.Separator), host, fileName))
	log.Debug(dir + name)
	os.MkdirAll(dir, 0744)
	viper := viper.New()
	viper.Set("prop1", "abcdefg")
	viper.Set("prop2", "hijklmn")
	viper.WriteConfigAs(dir + name)
	name = "config2.properties"
	viper.Set("prop2", "abcdefg")
	viper.Set("prop1", "hijklmn")
	viper.WriteConfigAs(dir + name)
}

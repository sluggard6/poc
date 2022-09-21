package service

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sluggard/poc/config"
	"github.com/sluggard/poc/util"
	"github.com/spf13/viper"
)

type FileService interface {
	LoadRemoteFile(localPath string, remotePath string) error
	SendRemoteFile(localPath string, remotePath string) error
	SaveFile(reader io.Reader, name string, host string) error
	LoadStringFile(name string, host string) (string, error)
	SearchFile(hosts []string, fileName string, prop string, propValue string, remote bool) []ReadProp
	UpdateProp(hosts []string, fileName string, prop string, propValue string, remote bool) error
}

var fileService = newFileImpl()

func GetFileService() FileService {
	return fileService
}

type fileImpl struct {
	cs   CommandService
	Tmp  string
	Root string
}

func newFileImpl() *fileImpl {
	path, _ := os.Getwd()
	return &fileImpl{GetCommandService(), ".tmp", path}
}

func (f *fileImpl) LoadRemoteFile(localPath string, remotePath string) error {
	cmd := "scp " + localPath + " " + remotePath
	if _, _, err := f.cs.Run(cmd); err != nil {
		return err
	}
	return nil
}

func (f *fileImpl) SendRemoteFile(localPath string, remotePath string) error {
	cmd := "scp " + remotePath + " " + localPath
	if _, _, err := f.cs.Run(cmd); err != nil {
		return err
	}
	return nil
}

func (f *fileImpl) SaveFile(reader io.Reader, name string, host string) error {
	tmpFile := f.newTmpFile()
	util.SaveAndSha(reader, tmpFile)
	// var fileName = fs.Root + string(filepath.Separator) + strings.Join(makeFilePath(hexString), string(filepath.Separator)) + filepath.Ext(name)
	var fileName = "." + host + string(filepath.Separator) + filepath.Ext(name)
	log.Debugf("store file : %s", fileName)
	dir, _ := filepath.Split(f.Root + string(filepath.Separator) + fileName)
	if err := os.MkdirAll(dir, 0744); err != nil {
		return err
	}
	os.Rename(tmpFile, f.Root+string(filepath.Separator)+fileName)
	return nil
}
func (f *fileImpl) newTmpFile() (filePath string) {
	name := util.UUID()
	filePath = f.Tmp + string(filepath.Separator) + name
	return
}

func (f *fileImpl) LoadStringFile(name string, host string) (string, error) {
	if bs, err := os.ReadFile("." + host + string(filepath.Separator) + name); err != nil {
		return "", err
	} else {
		return string(bs), nil
	}
}

func (f *fileImpl) ReadProp(path string) *viper.Viper {
	prop := viper.New()
	prop.SetConfigFile(path)
	prop.SetConfigType("properties")
	prop.ReadInConfig()
	return prop
}

func (f *fileImpl) SearchFile(hosts []string, fileName string, prop string, propValue string, remote bool) []ReadProp {
	if remote {
		return f.searchRemoteFile(hosts, fileName, prop, propValue)
	} else {
		return searchLocalFile(hosts, fileName, prop, propValue)
	}
}

type ReadProp struct {
	Host  string
	Prop  string
	Value string
}

func searchLocalFile(hosts []string, fileName string, prop string, propValue string) []ReadProp {
	ret := make([]ReadProp, 0)
	home, err := getHome()
	if err != nil {
		panic(fmt.Errorf("Fatal error get home: %s \n", err))
	}
	for _, host := range hosts {
		viper := viper.New()
		log.Debugf("%s%s.poc%s.%s%s", home, string(filepath.Separator), string(filepath.Separator), host, fileName)
		// dir, name := filepath.Split(fmt.Sprintf("%s%s.poe%s.%s%s", home, string(filepath.Separator), string(filepath.Separator), host, fileName))
		viper.SetConfigFile(fmt.Sprintf("%s%s.poc%s.%s%s", home, string(filepath.Separator), string(filepath.Separator), host, fileName))
		// viper.SetConfigFile(name)
		// viper.SetConfigName("config")
		viper.SetConfigType("properties")
		// viper.AddConfigPath(dir)
		err := viper.ReadInConfig() // 查找并读取配置文件
		if err != nil {             // 处理读取配置文件的错误
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		log.Debug(viper.Get(prop))
		if strings.Contains(viper.GetString(prop), propValue) {
			ret = append(ret, ReadProp{Host: host, Prop: prop, Value: viper.GetString(prop)})
		}
	}
	return ret
}

func (f *fileImpl) searchRemoteFile(hosts []string, fileName string, prop string, propValue string) []ReadProp {
	home, err := getHome()
	if err != nil {
		panic(fmt.Errorf("Fatal error get home: %s \n", err))
	}
	for _, ip := range hosts {
		dir, name := filepath.Split(fmt.Sprintf("%s%s.poe%s.%s%s", home, string(filepath.Separator), string(filepath.Separator), ip, fileName))
		os.MkdirAll(dir, 0744)
		cmd := fmt.Sprintf("sshpaas -p %s scp -rf %s@%s:%s %s", config.GetConfig().DevicesInfo.Password, config.GetConfig().DevicesInfo.Username, ip, fileName, dir+name)
		f.cs.Run(cmd)
	}
	return searchLocalFile(hosts, fileName, prop, propValue)
}

func (f *fileImpl) UpdateProp(hosts []string, fileName string, prop string, propValue string, remote bool) error {
	if remote {
		return f.updateRemoteProp(hosts, fileName, prop, propValue)
	} else {
		return updateLocalProp(hosts, fileName, prop, propValue)
	}
}

func updateLocalProp(hosts []string, fileName string, prop string, propValue string) error {
	home, err := getHome()
	if err != nil {
		// panic(fmt.Errorf("Fatal error get home: %s \n", err))
		return err
	}
	for _, host := range hosts {
		viper := viper.New()
		log.Debugf("%s%s.poc%s.%s%s", home, string(filepath.Separator), string(filepath.Separator), host, fileName)
		// dir, name := filepath.Split(fmt.Sprintf("%s%s.poe%s.%s%s", home, string(filepath.Separator), string(filepath.Separator), host, fileName))
		viper.SetConfigFile(fmt.Sprintf("%s%s.poc%s.%s%s", home, string(filepath.Separator), string(filepath.Separator), host, fileName))
		// viper.SetConfigFile(name)
		// viper.SetConfigName("config")
		viper.SetConfigType("properties")
		// viper.AddConfigPath(dir)
		err := viper.ReadInConfig() // 查找并读取配置文件
		if err != nil {             // 处理读取配置文件的错误
			// panic(fmt.Errorf("Fatal error config file: %s \n", err))
			return err
		}
		viper.Set(prop, propValue)
		err = viper.WriteConfig()
		if err != nil { // 处理写入配置文件的错误
			// panic(fmt.Errorf("Fatal error config file: %s \n", err))
			return err
		}
	}
	return nil
}

func (f *fileImpl) updateRemoteProp(hosts []string, fileName string, prop string, propValue string) error {
	updateLocalProp(hosts, fileName, prop, propValue)
	home, err := getHome()
	if err != nil {
		// panic(fmt.Errorf("Fatal error get home: %s \n", err))
		return err
	}
	for _, ip := range hosts {
		filePaht := fmt.Sprintf("%s%s.poe%s.%s%s", home, string(filepath.Separator), string(filepath.Separator), ip, fileName)
		cmd := fmt.Sprintf("sshpaas -p %s scp -rf %s %s@%s:%s", config.GetConfig().DevicesInfo.Password, filePaht, config.GetConfig().DevicesInfo.Username, ip, fileName)
		f.cs.Run(cmd)
	}
	return nil
}

func getHome() (string, error) {
	user, err := user.Current()
	if err == nil {
		return user.HomeDir, nil
	} else {
		return "", err
	}
}

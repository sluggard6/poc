package service

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/sluggard/poc/util"
	"github.com/spf13/viper"
)

type FileService interface {
	LoadRemoteFile(localPath string, remotePath string) error
	SendRemoteFile(localPath string, remotePath string) error
	SaveFile(reader io.Reader, name string, host string) error
	LoadStringFile(name string, host string) (string, error)
	SearchFile(hosts []string, fileName string, prop string, propValue string, remote bool) []ReadProp
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
	logrus.Debugf("store file : %s", fileName)
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
		return searchRemoteFile(hosts, fileName, prop, propValue)
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
	user, err := user.Current()
	home := ""
	ret := make([]ReadProp, 0)
	if err == nil {
		home = user.HomeDir
	}
	for _, host := range hosts {
		viper := viper.New()
		log.Debugf("%s%s.poc%s.%s%s", home, string(filepath.Separator), string(filepath.Separator), host, fileName)
		dir, name := filepath.Split(fmt.Sprintf("%s%s.poe%s.%s%s", home, string(filepath.Separator), string(filepath.Separator), host, fileName))
		viper.SetConfigFile(name)
		// viper.SetConfigName("config")
		viper.SetConfigType("ini")
		viper.AddConfigPath(dir)
		err := viper.ReadInConfig() // 查找并读取配置文件
		if err != nil {             // 处理读取配置文件的错误
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		for _, key := range viper.AllKeys() {
			log.Debug(key)
		}
		if strings.Contains(viper.GetString(prop), propValue) {
			ret = append(ret, ReadProp{Host: host, Prop: prop, Value: viper.GetString(prop)})
		}
	}
	return ret
}

func searchRemoteFile(hosts []string, fileName string, prop string, propValue string) []ReadProp {
	return make([]ReadProp, 0)
}

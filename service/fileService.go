package service

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/sluggard/poc/util"
)

type FileService interface {
	LoadRemoteFile(localPath string, remotePath string) error
	SendRemoteFile(localPath string, remotePath string) error
	SaveFile(reader io.Reader, name string, host string) error
	LoadStringFile(name string, host string) (string, error)
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
	f.cs.Run(cmd)
	return nil
}

func (f *fileImpl) SendRemoteFile(localPath string, remotePath string) error {
	cmd := "scp " + remotePath + " " + localPath
	f.cs.Run(cmd)
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

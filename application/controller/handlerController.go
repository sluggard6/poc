package controller

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"

	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"github.com/sluggard/poc/config"
	"github.com/sluggard/poc/service"
	"github.com/sluggard/poc/util"
)

type HandlerController struct {
	commandService service.CommandService
	fileService    service.FileService
}

func NewHandlerController() *HandlerController {
	return &HandlerController{service.GetCommandService(), service.GetFileService()}
}

func (c *HandlerController) PostCommand(ctx iris.Context) *HttpResult {
	cmd := "ls"
	if stdout, stderr, err := c.commandService.Run(cmd); err == nil {
		if len(stdout) > 0 {
			return Success(stdout)
		} else if len(stderr) > 0 {
			return Success(stderr)
		}
		return Success(stdout)
	} else {
		return FailedMessage(err.Error())
	}
}

func (c *HandlerController) GetScan(ctx iris.Context) *HttpResult {
	addrs, err := net.InterfaceAddrs()
	hosts := make([]string, 0)
	if err != nil {
		return FailedMessage(err.Error())
	}
	//获取本机地址
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				hosts = append(hosts, c.findIps(ipnet)...)
			}
		}
		// if len(hosts) > 0 {
		// 	util.Hosts = hosts
		// }
		log.Debug(hosts)
	}

	return Success(util.Hosts)
}

func (c *HandlerController) findIps(ipnet *net.IPNet) []string {
	cmd := "nmap -sP " + ipnet.IP.String() + "/24"
	// cmd = "nmap -sP " + "192.168.2.183" + "/24"
	log.Debug(cmd)
	if stdOut, _, err := c.commandService.Run(cmd); err == nil {
		log.Debug(stdOut)
		// util.Hosts = scanIp(stdOut, ipnet.IP.String())
		return scanIp(stdOut, ipnet.IP.String())
	} else {
		return make([]string, 0)
	}
}

func (c *HandlerController) GetHosts(ctx iris.Context) *HttpResult {
	return Success(util.Hosts)
}

func (c *HandlerController) GetFile(ctx iris.Context) *HttpResult {
	host := ctx.Params().Get("host")
	wd, _ := os.Getwd()
	localPath := wd + string(filepath.Separator) + "." + host + string(filepath.Separator) + util.FileName
	remotePath := fmt.Sprintf("%s@%s:%s", config.GetConfig().DevicesInfo.Username, host, util.ConfigFilePath)
	dir, _ := filepath.Split(localPath)
	if err := os.MkdirAll(dir, 0744); err != nil {
		return FailedMessage(err.Error())
	}
	if err := c.fileService.LoadRemoteFile(localPath, remotePath); err != nil {
		return FailedMessage(err.Error())
	}
	if fileString, err := c.fileService.LoadStringFile(util.FileName, host); err != nil {
		return FailedMessage(err.Error())
	} else {
		return Success(fileString)
	}
}

func scanIp(input string, localIp string) []string {
	reg1 := regexp.MustCompile(`((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`)
	s := make([]string, 0)
	for _, ss := range reg1.FindAllStringSubmatch(input, -1) {
		if localIp != ss[0] {
			s = append(s, ss[0])
		}
		// for j, s := range ss {
		// 	log.Debugf("s[%d][%d]:%s", i, j, s)
		// }
	}
	return s
}

func (c *HandlerController) PostSearch(ctx iris.Context) *HttpResult {
	// editLibraryForm := struct {
	// 	Id   uint   `json:"id"`
	// 	Name string `json:"name"`
	// }{}
	searchForm := struct {
		FileName  string `json:"fileName"`
		Prop      string `json:"prop"`
		PropValue string `json:"propValue"`
	}{}
	if err := ctx.ReadJSON(&searchForm); err != nil {
		return FailedCodeMessage(PARAM_ERROR, err.Error())
	}
	res := c.fileService.SearchFile(
		util.Filter(util.Hosts, func(h util.Host) bool {
			return h.State == 1
		}, func(h util.Host) string {
			return h.Ip
		}),
		searchForm.FileName, searchForm.Prop, searchForm.PropValue, config.GetConfig().RunConfig.Remote)
	return Success(res)
}

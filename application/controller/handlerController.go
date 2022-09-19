package controller

import (
	"net"
	"regexp"

	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"github.com/sluggard/poc/service"
	"github.com/sluggard/poc/util"
)

type HandlerController struct {
	commandService service.CommandService
}

func NewHandlerController() *HandlerController {
	return &HandlerController{service.GetCommandService()}
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
	var ipNet *net.IPNet
	if err != nil {
		return FailedMessage(err.Error())
	}
	//获取本机地址
	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipNet = ipnet
			}
		}
	}
	cmd := "nmap -sP " + ipNet.IP.String() + "/24"
	// cmd = "nmap -sP " + "192.168.2.183" + "/24"
	log.Debug(cmd)
	if stdOut, _, err := c.commandService.Run(cmd); err == nil {
		log.Debug(stdOut)
		util.Hosts = scanIp(stdOut, ipNet.IP.String())
		log.Debug(util.Hosts)
	}
	return Success(util.Hosts)
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

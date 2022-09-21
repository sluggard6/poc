//go:generate statik -src=./assets/dist
//go:generate go fmt statik/statik.go

package application

import (
	stdContext "context"
	"fmt"
	"time"

	"github.com/rakyll/statik/fs"
	"github.com/robfig/cron/v3"
	_ "github.com/sluggard/poc/statik" // TODO: Replace with the absolute import path

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	log "github.com/sirupsen/logrus"

	"github.com/sluggard/poc/application/controller"
	"github.com/sluggard/poc/config"
	"github.com/sluggard/poc/util"
)

// HttpServer
type HttpServer struct {
	Config  config.Config
	App     *iris.Application
	crontab *cron.Cron
	// Store  store.Store
	Status bool
}

var (
	gSessionId    = "GSESSIONID"
	sess          = sessions.New(sessions.Config{Cookie: gSessionId})
	ignoreAuthUrl = []string{
		"/test/ping",
		"/user/login",
		"/user/register",
		"/app",
	}
)

func NewServer(config config.Config) *HttpServer {
	app := iris.New()

	httpServer := &HttpServer{
		Config: config,
		App:    app,
		Status: false,
	}
	httpServer._Init()
	return httpServer
}

// Start
func (s *HttpServer) Start() error {
	if err := s.App.Run(
		// iris.Addr(fmt.Sprintf("%s:%d", libs.Config.Host, libs.Config.Port)),
		iris.Addr(fmt.Sprintf("%s:%d", s.Config.Server.Host, s.Config.Server.Port)),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(time.RFC3339),
	); err != nil {
		return err
	}
	s.Status = true
	return nil
}

// Start close the server at 3-6 seconds
func (s *HttpServer) Stop() {
	go func() {
		time.Sleep(3 * time.Second)
		ctx, cancel := stdContext.WithTimeout(stdContext.TODO(), 3*time.Second)
		defer cancel()
		s.crontab.Stop()
		s.App.Shutdown(ctx)
		s.Status = false
	}()
}

func (s *HttpServer) _Init() error {
	s.RouteInit()
	util.Hosts = make([]util.Host, 0)
	for _, ip := range config.GetConfig().DevicesInfo.Hosts {
		util.Hosts = append(util.Hosts, util.Host{Ip: ip, State: 1})
		config.MakeDemoFile(ip)
	}
	if config.GetConfig().RunConfig.Remote {
		go controller.NewHandlerController().GetScan(nil)
		// 新建一个定时任务对象
		// 根据cron表达式进行时间调度，cron可以精确到秒，大部分表达式格式也是从秒开始。
		s.crontab = cron.New(cron.WithSeconds()) //精确到秒
		// 定义定时器调用的任务函数
		//定时任务
		spec := "*/10 * * * *" //cron表达式，每五秒一次
		task := func() {
			//fmt.Println("hello world", time.Now())
			go controller.NewHandlerController().GetScan(nil)
		}
		s.crontab.AddFunc(spec, task)
		s.crontab.Start()
	}
	return nil
}

// RouteInit 初始化路由
func (s *HttpServer) RouteInit() {

	app := s.App
	app.Options("/*", controller.Cors)
	// app.Party("/*", controller.Cors).AllowMethods(iris.MethodOptions)
	app.UseGlobal(controller.Cors)
	statikFS, err := fs.New()
	if err == nil {
		// app.Handle()
		app.HandleDir(s.Config.Server.ContextPath+"/app", statikFS)
	} else {
		fmt.Printf("err: %v\n", err)
	}
	//app.Use(AuthRequired, sess.Handler())
	// app.Use(sess.Handler())
	// statikFS, err := fs.New()
	// if err == nil {
	// 	app.HandleDir(s.Config.Server.ContextPath+"/app", statikFS)
	// } else {
	// 	fmt.Printf("err: %v\n", err)
	// }
	mvc.New(app.Party(s.Config.Server.ContextPath + "/")).Handle(controller.NewHandlerController())
	for _, route := range app.APIBuilder.GetRoutes() {
		log.Info(route)
	}

}

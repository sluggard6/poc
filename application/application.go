//go:generate statik -src=./assets/dist
//go:generate go fmt statik/statik.go

package application

import (
	stdContext "context"
	"fmt"
	"strings"
	"time"

	//_ "github.com/sluggard/myfile/statik" // TODO: Replace with the absolute import path

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	log "github.com/sirupsen/logrus"

	"github.com/sluggard/poc/application/controller"
	"github.com/sluggard/poc/config"
)

// HttpServer
type HttpServer struct {
	Config config.Config
	App    *iris.Application
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
		s.App.Shutdown(ctx)
		s.Status = false
	}()
}

func (s *HttpServer) _Init() error {
	s.RouteInit()
	return nil
}

// AuthRequired 登录验证
func AuthRequired(ctx iris.Context) {
	// iris.CookieSameSite = iris.SameSiteNoneMode
	log.Debug(ctx.Request().Method, "  ", ctx.Request().RequestURI)
	// log.Debug(ctx.RequestPath(false))
	// session := sess.Start(ctx, iris.CookieSameSite(iris.SameSiteNoneMode))
	// ctx.addCookieOption
	session := sess.Start(ctx)
	// log.Debug(session)
	// log.Debug(sess.GetCookieOptions())
	path := config.GetConfig().Server.ContextPath
	//被忽略的url直接通过
	for _, v := range ignoreAuthUrl {
		if strings.HasPrefix(ctx.RequestPath(false), path+v) {
			// if path+v == ctx.RequestPath(false) {
			ctx.Next()
			return
		}
	}
	if auth, _ := session.GetBoolean("authenticated"); !auth {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}
	ctx.Next()
}

// RouteInit 初始化路由
func (s *HttpServer) RouteInit() {

	app := s.App
	app.Options("/*", controller.Cors)
	// app.Party("/*", controller.Cors).AllowMethods(iris.MethodOptions)
	app.UseGlobal(controller.Cors)
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
	go controller.NewHandlerController().GetScan(nil)
}

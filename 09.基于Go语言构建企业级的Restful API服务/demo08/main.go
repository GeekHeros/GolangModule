package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"ApiServer/demo08/router"
	"errors"
	"github.com/spf13/pflag"
	"ApiServer/demo08/config"
	"github.com/spf13/viper"
	"github.com/lexkong/log"
	"ApiServer/demo08/model"
	"ApiServer/demo08/router/middleware"
)

var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path.")
)

func main() {
	// 解析flag
	pflag.Parse()

	// 初始化配置
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// init db
	model.DB.Init()
	defer model.DB.Close()

	//我修改配置文件之后提示这个：errno: open E:\Go\GOPATH\src\gitee.wang\a\restfulapi\conf\config.yaml: The process cannot access the file because it is being used by another process.
	//可能是因为你用IDEA这种工具来编辑这个config.yaml，试下用vim来修改
	//for{
	//	fmt.Println(viper.GetString("runmode"))
	//	time.Sleep(4*time.Second)
	//}

	// 设定Gin的模式
	gin.SetMode(viper.GetString("runmode"))

	// 添加日志测试功能:日志转传
	//for {
	//	log.Info("111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")
	//	time.Sleep(100 * time.Millisecond)
	//}

	// Create the Gin engine.
	g := gin.New()

	middlewares := []gin.HandlerFunc{
		middleware.RequestId(),
		middleware.Logging(),
	}

	// Routes.
	router.Load(
		// Cores.
		g,

		// Middlwares.
		middlewares...,
	)

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()

	log.Info("Start to listening the incoming requests on http address: " + viper.GetString("addr"))
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}

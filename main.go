/*
 @Title  请填写文件名称（需要改）
 @Description  请填写文件描述（需要改）
 @Author  Leo  2020/4/28 3:56 下午
 @Update  Leo  2020/4/28 3:56 下午
*/
package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	ini "gopkg.in/ini.v1"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	confInst *ini.File
	sigs = make(chan os.Signal)
	done = make(chan bool)
)

func main() {

	confInst,_ = ini.Load("./c.ini")

	flag.Parse()
	err:=flag.Set("log_dir", "./logs/")
	if err!=nil {
		panic(err)
	}

	go serve()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // ctrl+c, kill, kill -2,
	// 注意 kill -9 属于强制结束，不会被处理也无法被任何程序处理
	// https://blog.csdn.net/IT_flying625/article/details/89554946 Linux signal 信号列表（sigint,sigtstp..)
	// https://blog.csdn.net/e421083458/article/details/52712009 golang如何优雅处理退出

	go sigAwaiter()


	<-done
}

//func testSiged() {
//
//}

func sigAwaiter() {
	sig := <-sigs
	fmt.Println(fmt.Sprintf("recv signal %s", sig.String()))

	glog.Flush() // kill , ctrl+c 可以， kill -9 不行
	done<-true
}

func serve() {
	appName := confInst.Section("app").Key("name").String()

	//log := golog.New("testlog")

	a := uuid.NewV1()
	e := echo.New()

	e.Any("/", func(c echo.Context) error {

		name := c.QueryParams().Get("name")

		//log.Debugf("recv request %s", name)
		glog.Infof("recv request %s", name)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"app_name": appName ,
			"data":"hello "+a.String() ,
		})

	})
	err := e.Start(":8080")

	if err!=nil {
		panic(err)
	}
}
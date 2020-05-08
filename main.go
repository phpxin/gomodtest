/*
 @Title  请填写文件名称（需要改）
 @Description  请填写文件描述（需要改）
 @Author  Leo  2020/4/28 3:56 下午
 @Update  Leo  2020/4/28 3:56 下午
*/
package main

import (
	"fmt"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	ini "gopkg.in/ini.v1"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	confInst *ini.File
	sigs = make(chan os.Signal)
	done = make(chan bool)
	logger *zap.Logger
)

func main() {
	defer fmt.Println("will exit")


	initZapLogger()

	logger.Info("server started")


	confInst,_ = ini.Load("./c.ini")

	//flag.Parse()
	//err =flag.Set("log_dir", "./logs/")
	//if err!=nil {
	//	panic(err)
	//}

	go serve()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // ctrl+c, kill, kill -2,
	// 注意 kill -9 属于强制结束，不会被处理也无法被任何程序处理
	// https://blog.csdn.net/IT_flying625/article/details/89554946 Linux signal 信号列表（sigint,sigtstp..)
	// https://blog.csdn.net/e421083458/article/details/52712009 golang如何优雅处理退出

	go sigAwaiter()


	<-done
}

func initZapLogger() {

	hook := lumberjack.Logger{
		Filename:   "/Users/leo/Downloads/logs/godemo/mygomodtest-1.log", // 日志文件路径
		MaxSize:    128,                      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                       // 日志文件最多保存多少个备份
		MaxAge:     7,                        // 文件最多保存多少天
		Compress:   true,                     // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel,                                                                     // 日志级别
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("serviceName", "serviceName"))
	// 构造日志
	logger = zap.New(core, caller, development, filed)

	logger.Info("log 初始化成功")
}

func sigAwaiter() {
	sig := <-sigs
	fmt.Println(fmt.Sprintf("recv signal %s", sig.String()))

	//glog.Flush() // kill , ctrl+c 可以， kill -9 不行

	var err error

	err = logger.Sync()
	if err!=nil {
		fmt.Println("Err ", err)
	}

	done<-true
}

func serve2() {
	defer PrintPanicStackError() // recover from panic

	ch := make(chan string, 100)

	go serveReader(ch)

	serveWriter(ch)

	fmt.Println("server writer exited") // 当 serveWriter 中触发了 panic 这里是不会被执行的
}

func serveWriter(ch chan<- string) {

	for {
		str := ""
		fmt.Println("please input some information:")
		_,e := fmt.Scanf("%s", &str)
		if e!=nil {
			fmt.Println("ERR ", e.Error())
			continue
			//panic(e)
		}

		ch <- str
	}

}

func serveReader(ch  <-chan string) {
	defer fmt.Println("after catch panic [serveReader]")
	defer PrintPanicStackError() // recover from panic
	defer fmt.Println("before catch panic [serveReader]")

	for {
		select {
		case str := <-ch :
			if str=="panictest" {
				panictest()
			}else{
				fmt.Println("input is ", str)
			}
		}
	}
}

func panictest() {
	panic("test panic")
}

func serve() {
	appName := confInst.Section("app").Key("name").String()

	a := uuid.NewV1()
	e := echo.New()

	e.Any("/", func(c echo.Context) error {

		//name := c.QueryParams().Get("name")

		logger.Info("http request", zap.String("url", c.Path()) ,
			zap.String("query_string", c.QueryString()))

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
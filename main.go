/*
 @Title  请填写文件名称（需要改）
 @Description  请填写文件描述（需要改）
 @Author  Leo  2020/4/28 3:56 下午
 @Update  Leo  2020/4/28 3:56 下午
*/
package main

import (
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	ini "gopkg.in/ini.v1"
	"net/http"
)

func main() {

	confInst,_ := ini.Load("./c.ini")
	appName := confInst.Section("app").Key("name").String()

	a := uuid.NewV1()
	e := echo.New()

	e.Any("/", func(c echo.Context) error {
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
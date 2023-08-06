package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
)

// 定义中间件
func myHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("usersession", "userid-1")
		ctx.Next() //放行

		//ctx.Abort() //阻止
	}
}

func main() {
	//创建服务
	ginServer := gin.Default()

	//注册中间件
	ginServer.Use(myHandler())

	//设置网页图标
	ginServer.Use(favicon.New("./OIG.jpg"))

	//加载静态页面
	ginServer.LoadHTMLGlob("templates/*")

	//加载静态文件 相对路径 绝对路径
	ginServer.Static("/static", "./static")
	//相应页面给前端
	ginServer.GET("/index", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{
			"msg": "后台消息",
		})
	})

	// /user/info?userId=1&userName=xingzhou
	ginServer.GET("/user/info", myHandler(), func(ctx *gin.Context) {
		usersession := ctx.MustGet("usersession").(string)
		log.Println("================>", usersession)
		userId := ctx.Query("userId")
		userName := ctx.Query("userName")
		ctx.JSON(200, gin.H{
			"userId":   userId,
			"userName": userName,
		})

	})

	// /user/info/1/xingzhou
	ginServer.GET("/user/info/:userId/:userName", func(ctx *gin.Context) {
		userId := ctx.Param("userId")
		userName := ctx.Param("userName")
		ctx.JSON(200, gin.H{
			"userId":   userId,
			"userName": userName,
		})
	})

	//前端返回数据给后端
	ginServer.POST("/json", func(ctx *gin.Context) {
		data, _ := ctx.GetRawData()

		var m map[string]interface{}
		//包装为json []byte
		_ = json.Unmarshal(data, &m)
		ctx.JSON(200, m)
	})

	//获取表单数据
	ginServer.POST("/user/add", func(ctx *gin.Context) {
		userName := ctx.PostForm("userName")
		password := ctx.PostForm("password")

		ctx.JSON(200, gin.H{
			"smg":      "ok",
			"userName": userName,
			"password": password,
		})
	})

	//文件上传
	ginServer.GET("/files", func(ctx *gin.Context) {
		ctx.HTML(200, "files.html", nil)
	})
	//	处理文件流请求
	ginServer.POST("/upload", func(ctx *gin.Context) {

		//从请求中读取文件
		f, err := ctx.FormFile("f1")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
		} else {
			//将文件保存至本地
			ctx.HTML(200, "ok.html", nil)
			dst := fmt.Sprintf("./files/%s", f.Filename)

			ctx.SaveUploadedFile(f, dst)
		}
	})

	//路由
	ginServer.GET("/test", func(ctx *gin.Context) {
		//重定向 301
		ctx.Redirect(301, "https://www.bilibili.com")
	})
	//404
	ginServer.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(404, "404.html", nil)
	})

	//路由组
	userGroup := ginServer.Group("/User")
	{
		userGroup.GET("/add")
		userGroup.POST("login")
		userGroup.POST("logout")
	}
	orderGroup := ginServer.Group("/order")
	{
		orderGroup.GET("add")
		orderGroup.DELETE("/del")
	}
	ginServer.Run(":8082")
}

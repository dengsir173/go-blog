package router

import (
	"blog/controller"
	_ "embed"
	_ "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
)

func Start() {

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")

	//router.GET("/index",controller.ListUser)
	router.POST("/register", controller.Register)
	router.GET("/register", controller.GoRegister)
	router.GET("/", controller.Index)
	router.GET("/login", controller.GoLogin)
	router.POST("/login", controller.Login)
	//操作博客
	router.GET("/blog_index", controller.GetBlogIndex)
	router.POST("/blog", controller.AddBlog)
	router.GET("/blog", controller.GoAddBlog)
	router.GET("/blogdetail", controller.BlogDetail)
	//操作背景图
	router.GET("/gouploadimage", controller.GoUploadImage)
	router.POST("/uploadimage", controller.UploadImage)
	router.Run(":80")
}

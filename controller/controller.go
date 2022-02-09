package controller

import (
	"blog/dao"
	"blog/model"
	_ "embed"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/russross/blackfriday"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

func Register(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	person := model.Person{
		Username: username,
		Password: password,
	}
	dao.Mgr.Register(&person)
	c.Redirect(301, "/")
}

func GoRegister(c *gin.Context) {
	c.HTML(200, "register.html", nil)
}

func ListUser(c *gin.Context) {
	c.HTML(200, "userList.html", nil)
}

func Index(c *gin.Context) {
	//此处加入redis缓存
	conn, err := redis.Dial("tcp", "124.71.14.55:6379")
	if err != nil {
		fmt.Println("链接redis失败", err)
		return
	}
	defer conn.Close()
	b, _ := redis.Bool(conn.Do("EXISTS", "imageUrl"))
	if b {
		url, _ := redis.String(conn.Do("GET", "imageUrl"))
		c.HTML(200, "index.html", url)
	} else {
		lastImage := dao.Mgr.GetLastImage()
		_, err := conn.Do("SET", "imageUrl", lastImage.Url)
		if err != nil {
			c.String(400, fmt.Sprintf("redis set error", err.Error()))
			return
		}
		c.HTML(200, "index.html", lastImage.Url)
	}
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	fmt.Println(username)
	u := dao.Mgr.Login(username)

	if u.Username == "" {
		c.HTML(200, "login.html", "用户名不存在！")
		fmt.Println("用户名不存在")
	} else {
		if u.Password != password {
			fmt.Println("密码错误")
			c.HTML(200, "login.html", "密码错误")
		} else {
			fmt.Println("登录成功")
			c.Redirect(301, "/")
		}
	}
}

func GoLogin(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

//操作博客
func GetBlogIndex(c *gin.Context) {
	blogs := dao.Mgr.GetAllBlog()
	c.HTML(200, "blogIndex.html", blogs)
}
func AddBlog(c *gin.Context) {
	title := c.PostForm("title")
	tag := c.PostForm("tag")
	content := c.PostForm("content")

	blog := model.Blog{
		Title:   title,
		Tag:     tag,
		Content: content,
	}
	dao.Mgr.AddBlog(&blog)
	c.Redirect(302, "/blog_index")
}
func GoAddBlog(c *gin.Context) {
	c.HTML(200, "blog.html", nil)
}

func BlogDetail(c *gin.Context) {
	s := c.Query("bid")
	bid, _ := strconv.Atoi(s)
	b := dao.Mgr.GetBlog(bid)
	content := blackfriday.MarkdownCommon([]byte(b.Content))
	c.HTML(200, "detail.html", gin.H{
		"Title":   b.Title,
		"Content": template.HTML(content),
	})
}

//图片上传
func GoUploadImage(c *gin.Context) {
	c.HTML(200, "UploadImage.html", nil)
}

var filename = uuid.NewString() + ".jpg"

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}
	pre := "./assets/img/"
	if err := c.SaveUploadedFile(file, pre+filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("图片上传失败：%s", err.Error()))
		return
	}
	url := Oss(filename)
	image := model.Image{
		Url: url,
	}
	dao.Mgr.AddImage(&image)
	conn, err := redis.Dial("tcp", "124.71.14.55:6379")
	if err != nil {
		fmt.Println("链接redis失败", err)
		return
	}
	defer conn.Close()
	_, err = conn.Do("DEL", "imageUrl")
	if err != nil {
		fmt.Println("redis delelte failed:", err)
	}
	c.Redirect(301, "/")
	//c.HTML(200, "index.html", url)

}

func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

//oss对象存储
func Oss(filename string) string {
	// Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := "oss-cn-shenzhen.aliyuncs.com"
	// 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
	accessKeyId := "LTAI5t6BSawcVdzZXK65Dvn7"
	accessKeySecret := "HwzSrFLi0iPadPK1PvNbQHbxEi4Y6C"
	bucketName := "geleigo"
	// <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	objectName := "blog/" + filename
	// <yourLocalFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
	localFileName := "./assets/img" + "/" + filename
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		handleError(err)
	}
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		handleError(err)
	}
	img := "http://" + bucketName + "." + endpoint + "/" + objectName
	return img
}

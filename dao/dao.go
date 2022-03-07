package dao

import (
	"blog/model"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Manager interface {
	Register(user *model.Person)
	Login(username string) model.Person

	//博客操作
	AddBlog(post *model.Blog)
	GetAllBlog() []model.Blog
	GetBlog(bid int) model.Blog

	//图片操作
	AddImage(image *model.Image)
	GetLastImage() model.Image
}

type manager struct {
	db *gorm.DB
}

var Mgr Manager

func init() {
	db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("database connect error")
		return
	}
	Mgr = &manager{db: db}
	db.AutoMigrate(new(model.Person))
	db.AutoMigrate(new(model.Blog))
	db.AutoMigrate(new(model.Image))
}

func (mgr *manager) Register(person *model.Person) {
	mgr.db.Create(person)
}

func (mgr *manager) Login(username string) model.Person {
	var person model.Person
	mgr.db.Where("username=?", username).First(&person)
	return person
}

//博客操作
func (mgr *manager) AddBlog(blog *model.Blog) {
	mgr.db.Create(blog)
}
func (mgr *manager) GetAllBlog() []model.Blog {
	var blogs = make([]model.Blog, 10)
	mgr.db.Find(&blogs)
	return blogs
}

func (mgr *manager) GetBlog(bid int) model.Blog {
	defer func() {
		err := recover() //内置函数，能捕获到异常
		if err != nil {
			fmt.Println("出错了")
		}
	}()

	var blog model.Blog
	mgr.db.First(&blog, bid)
	return blog

}

//图片操作
func (mgr *manager) AddImage(image *model.Image) {
	mgr.db.Create(image)
}
func (mgr *manager) GetLastImage() model.Image {
	var image model.Image
	mgr.db.Last(&image)
	return image
}

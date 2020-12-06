package user

import (
	"fmt"
	"github.com/enoperm/internet-services-2020/util"
	"gorm.io/gorm"
	"net/http"
	. "server/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)
import "log"

const (
	sessionUserId = "user-id"
)

type UserManager struct {
	db *gorm.DB
}

func AttachUserEndpoints(router gin.IRouter, db *gorm.DB) *UserManager {
	um := UserManager{
		db: db,
	}
	
	err := um.db.AutoMigrate(&User{})
	if err != nil { panic(err) }
	
	attachRegisterEndpoints(&um, router)
	attachLoginEndpoints(&um, router)
	
	return &um
}

func (userManager UserManager) GetCurrentUserFromSession(c *gin.Context) *User {
	var user User

	sess := sessions.Default(c)
	userId := sess.Get(sessionUserId)
	if userId != nil {
		uid := userId.(uint)
		result := userManager.db.First(&user, uid)
		if result.Error != nil {
			return nil
		}
		return &user
	}

	return nil
}

func (userManager UserManager) SetCurrentUser(c *gin.Context, user *User) {
	sess := sessions.Default(c)

	switch {
	case user != nil:
		sess.Set(sessionUserId, user.ID)
	default:
		sess.Delete(sessionUserId)
	}

	err := sess.Save()
	if err != nil {
		panic(err)
	}
}

func attachRegisterEndpoints(um *UserManager, router gin.IRouter) {
	router.GET("/register", um.renderRegister)
	router.POST("/register", um.postRegister)
}

func attachLoginEndpoints(um *UserManager, router gin.IRouter) {
	router.GET("/login", um.renderLogin)
	router.POST("/login", um.postLogin)
}

func (um UserManager) renderLogin(c *gin.Context) {
	util.HtmlWithContext(c, http.StatusOK, "login", gin.H{})
}

func (um UserManager) renderRegister(c *gin.Context) {
	util.HtmlWithContext(c, http.StatusOK, "register", gin.H{})
}

func (um UserManager) postRegister(c *gin.Context) {
	var regReq struct {
		Name string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}
	
	fmt.Printf("PR: %#v\n", c.Params)

	err := c.ShouldBind(&regReq)
	if err != nil { panic(err) }

	user := User{
		Name: regReq.Name,
		PasswordHash: HashPassword(regReq.Password),
	}

	tx := um.db.Create(&user)
	if tx.Error != nil {
		util.HtmlWithContext(c, http.StatusConflict, "register", gin.H{
			"message": "username already taken",
		})
		return
	}
	
	c.Redirect(http.StatusSeeOther, "/")
}

func (um UserManager) postLogin(c *gin.Context) {
	var credentials struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}
	
	err := c.ShouldBind(&credentials)
	if err != nil {
		um.renderLogin(c)
		return
	}

	var user User
	tx := um.db.Model(&user).Where("name = ?", credentials.Username).First(&user)
	if tx.Error != nil {
		panic(tx.Error)
	}
	
	err = CheckPasswordsMatch(credentials.Password, user.PasswordHash)
	if err != nil {
		log.Println("pwcheck:", err)
		um.renderLogin(c)
		return
	}
	
	um.SetCurrentUser(c, &user)
	
	c.Redirect(http.StatusSeeOther, "/")
}
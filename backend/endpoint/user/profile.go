package user

// TODO: Split file
import (
	"fmt"
	"github.com/enoperm/internet-services-2020/middleware"
	. "github.com/enoperm/internet-services-2020/model"
	"github.com/enoperm/internet-services-2020/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileManager struct {
	db *gorm.DB
}

func AttachProfileEnpoints(router gin.IRouter, db *gorm.DB) *ProfileManager {
	err := db.AutoMigrate(&Profile{})
	if err != nil {
		panic(err)
	}
	
	pm := ProfileManager{
		db: db,
	}

	router.GET("/profile", pm.getProfile)
	router.POST("/profile", pm.postProfile)

	return &pm
}

func (pm ProfileManager) renderProfile(c *gin.Context) {
	var profile Profile
	user := middleware.CurrentUser(c)
	pm.db.Model(&profile).Where("user_id = ?", user.ID).First(&profile)

	util.HtmlWithContext(c, http.StatusOK, "profile", gin.H{
		"profile": profile,
	})
}

func (pm ProfileManager) getProfile(c *gin.Context) {
	pm.renderProfile(c)
}

func (pm ProfileManager) postProfile(c *gin.Context) {
	var profile Profile
	user := middleware.CurrentUser(c)
	pm.db.Model(&profile).Where("user_id = ?", user.ID).First(&profile)
	
	fmt.Println(c.Params)
	
	err := c.ShouldBind(&profile)
	fmt.Println("bind:", err)
	profile.UserID = user.ID
	
	tx := pm.db.Save(&profile)
	if tx.Error != nil {
		panic(tx.Error)
	}
	
	c.Redirect(http.StatusSeeOther, "/auth/profile")
}

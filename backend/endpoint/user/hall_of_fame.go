package user

// TODO: Split file
import (
	. "github.com/enoperm/internet-services-2020/model"
	"github.com/enoperm/internet-services-2020/util"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HallOfFame struct {
	db *gorm.DB
}

func AttachHallOfFameEndpoints(router gin.IRouter, db *gorm.DB) *HallOfFame {
	err := db.AutoMigrate(&Profile{})
	if err != nil {
		panic(err)
	}
	
	hof := HallOfFame{
		db: db,
	}

	router.GET("/hall-of-fame", hof.getHall)

	return &hof
}

func (hof HallOfFame) renderHall(c *gin.Context) {
	var profiles []Profile

	tx := hof.db.Model(&Profile{});
	if tx.Error != nil {
		panic(tx.Error)
	}

	rows, _ := tx.Rows()
	for rows.Next() {
		var profile Profile
		_ = tx.ScanRows(rows, &profile)
		profiles = append(profiles, profile)
	}

	sort.Slice(profiles, func(i, j int) bool {
		left:= profiles[i]
		right := profiles[j]
		
		return left.LastSmokeTime().Before(right.LastSmokeTime())
	})

	type pubProf struct {
		Name string
		Days int
	}
	
	pubProfiles := make(map[int]pubProf, len(profiles))
	
	present := time.Now()
	days := func(t time.Time) int {
		return int(present.Sub(t).Hours() / 24.0)
	}

	for i, p := range profiles {
		var u User
		hof.db.Model(&User{}).Where("id = ?", p.UserID).First(&u)
		pubProfiles[i + 1] = pubProf{
			Name: u.Name,
			Days: days(p.LastSmokeTime()),
		}
	}

	util.HtmlWithContext(c, http.StatusOK, "hall-of-fame", gin.H{
		"Profiles": pubProfiles,
	})
}

func (hof HallOfFame) getHall(c *gin.Context) {
	hof.renderHall(c)
}


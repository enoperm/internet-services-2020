package main

import (
	"fmt"
	"github.com/enoperm/internet-services-2020/util"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"html/template"

	"net/http"

	"github.com/gin-gonic/nosurf"
	adapter "github.com/gwatts/gin-adapter"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"server/model"

	"server/endpoint/user"
	"server/middleware"
)

func main() {
	db, err := gorm.Open(sqlite.Open("./data/test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	cookieStore := cookie.NewStore([]byte("TODO-take-secret-from-elsewhere"))

	router := gin.Default()

	csrfHandler := func() gin.HandlerFunc {
		next, wrapper := adapter.New()
		nsHandler := nosurf.New(next)
		nsHandler.SetBaseCookie(http.Cookie{
			Path:     "/",
			HttpOnly: true,
		})
		nsHandler.SetFailureHandler(http.HandlerFunc((func(rw http.ResponseWriter, req *http.Request) {
			http.Error(rw, "failed to verify CSRF token", http.StatusBadRequest)
		})))

		return wrapper(nsHandler)
	}()

	router.Use(csrfHandler)
	router.Use(sessions.Sessions("login_state", cookieStore))

	router.Static("/static/", "./static/")

	router.HTMLRender = ginview.New(goview.Config{
		Root:         "views",
		Extension:    ".html.tmpl",
		Master:       "layouts/main",
		Funcs:        template.FuncMap{
			"humanTimespan": func(days uint32) string {
				years := days / 365
				days = days % 365
				months := days / 30
				days = days % 30
				
				var yearStr, monthStr, dayStr string
				
				if years > 0 {
					yearStr = fmt.Sprintf("%d years, ", years)
				}

				monthStr = fmt.Sprintf("%d months, ", months)
				dayStr = fmt.Sprintf("%d days", days)

				return fmt.Sprintf("%s%s%s", yearStr, monthStr, dayStr)
			},
			"benefitsAt": func(days uint32) string {
				switch {
				case days >= 60:
					return "Message 2"

				case days >= 30:
					return "Message 1"

				default:
					return "Message 0"
				}
			},
		},
		DisableCache: true,
	})

	userEndpoint := user.AttachUserEndpoints(router, db)
	userFromSession := func(c *gin.Context) *model.User {
		return userEndpoint.GetCurrentUserFromSession(c)
	}
	router.Use(middleware.WithUser(userFromSession))

	authorized := router.Group("/auth")
	{
		user.AttachProfileEnpoints(authorized, db)
	}

	// avoid import cycle so user endpoint package can defer tasks to middleware

	authorized.Use(middleware.AuthRequired())
	{
		// TODO: Profile, hall of fame
		authorized.POST("/logout", func(c *gin.Context) {
			userEndpoint.SetCurrentUser(c, nil)
			c.Redirect(http.StatusSeeOther, "/")
		})

		authorized.GET("/main", func(c *gin.Context) {
			var profile model.Profile
			u := middleware.CurrentUser(c)
			db.Model(&profile).Where("user_id = ?", u.ID).First(&profile)
			htmlContext := gin.H{}

			if profile.UserID == u.ID {
				htmlContext = gin.H{
					"Stats": profile.Stats(),
				}
			}

			util.HtmlWithContext(c, http.StatusOK, "main", htmlContext)
		})

		user.AttachHallOfFameEndpoints(authorized, db)
	}

	router.GET("/", func(c *gin.Context) {
		currentUser := middleware.CurrentUser(c)
		if currentUser != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/auth/main")
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, "/login")
	})

	http.ListenAndServe(":2000", router)
}

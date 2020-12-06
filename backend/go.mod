module github.com/enoperm/internet-services-2020

go 1.15

require (
	github.com/foolin/goview v0.3.0
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/gin-gonic/nosurf v0.0.0-20150415101651-45adcfcaf706
	github.com/gorilla/mux v1.8.0
	github.com/gwatts/gin-adapter v0.0.0-20170508204228-c44433c485ad
	github.com/mattn/go-sqlite3 v1.14.5
	golang.org/x/crypto v0.0.0-20201112155050-0c6587e931a9
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.20.8
	server v0.0.0-00010101000000-000000000000
)

replace server => ./

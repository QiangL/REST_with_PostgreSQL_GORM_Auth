package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chenyahui/gin-cache/persist"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var InfoLog = log.New(os.Stdout, "\x1b[31m[INF]\x1b[0m ", log.Ldate|log.Ltime)
var ErrorLog = log.New(os.Stderr, "\x1b[31m[ERR]\x1b[0m ", log.Ldate|log.Ltime)

var server = os.Getenv("DB_HOST")
var port = 1433
var user = os.Getenv("DB_USER_NAME")
var password = os.Getenv("DB_PASSWORD")
var database = os.Getenv("DB_NAME")

var DB *gorm.DB
var MemoryStore = persist.NewMemoryStore(3 * time.Hour)
var redisClient = redis.NewClient(&redis.Options{
	Network:     "tcp",
	Addr:        os.Getenv("REDIS_ENDPOINT"),
	ReadTimeout: 1 * time.Second,
	DialTimeout: 1 * time.Second,
})
var RedisStore = persist.NewRedisStore(redisClient)

//var RedisStore = persist.NewMemoryStore(3 * time.Hour)

func ConnectDB() {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	db, err := gorm.Open("mssql", connString)
	if err != nil {
		panic(err.Error())
	}
	db.AutoMigrate(&Ae86user{})
	//db.LogMode(true)
	DB = db
}

func BuildAuthKey(username string) string {
	return "auth-user-" + username
}

func BuildRateLimitKey(username string, date time.Time) string {
	return "rate-limit-" + username + date.Format("-2006-01-02")
}

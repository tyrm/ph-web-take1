package models

import (
	"database/sql"
	"math/rand"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/juju/loggo"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/rubenv/sql-migrate"
)

var DB *sql.DB
var logger *loggo.Logger

var cUsernameByID *cache.Cache

func CloseDB() {
	DB.Close()

	return
}

func Init(connectionString string) {
	newLogger := loggo.GetLogger("models")
	logger = &newLogger

	logger.Debugf("Connecting to Database")
	dbClient, err := sql.Open("postgres", connectionString)
	if err != nil {
		logger.Criticalf("Coud not connect to database: %s", err)
		panic(err)
	}
	DB = dbClient

	DB.SetMaxIdleConns(5)

	// Do Migration
	logger.Debugf("Loading Migrations")
	migrate.SetTable("web_migrations")
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("migrations","./migrations"),
	}

	logger.Debugf("Applying Migrations")
	n, err := migrate.Exec(DB, "postgres", migrations, migrate.Up)
	if n > 0 {
		logger.Infof("Applied %d migrations!\n", n)
	}
	if err != nil {
		logger.Criticalf("Coud not migrate database: %s", err)
		panic(err)
	}

	// init cache
	cUsernameByID = cache.New(5*time.Minute, 10*time.Minute)

	return
}

// Random string generator
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)
var randStringSrc = rand.NewSource(time.Now().UnixNano())
func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randStringSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randStringSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
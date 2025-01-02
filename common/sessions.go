package common

import (
	"context"
	"easywrite-service/db"
	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v8"
	"log"
)

var Sessions *redisstore.RedisStore

func InitSessions() {
	// New default RedisStore
	var err error
	Sessions, err = redisstore.NewRedisStore(context.Background(), db.Redis)
	if err != nil {
		log.Fatal("failed to create redis store: ", err)
	}

	Sessions.KeyPrefix("session_")
	Sessions.Options(sessions.Options{
		Path:   "/",
		MaxAge: 5 * 12 * 30 * 24 * 60 * 60,
	})
}

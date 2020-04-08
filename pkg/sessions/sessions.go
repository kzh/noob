package sessions

import (
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func Sessions() gin.HandlerFunc {
	log.Println("Connecting to Redis Sessions...")

	store, err := redis.NewStore(
		10,
		"tcp",
		"noob-sessions-master:6379",
		os.Getenv("SESSIONS_PASSWORD"),
		[]byte("NOOB_SESSION_SECRET"),
		//[]byte(os.Getenv("SESSION_SECRET")),
	)
	if err != nil {
		panic(err)
	}

	log.Println("Connected to Redis Sessions.")
	return sessions.Sessions("noob", store)
}

type NoobSession struct {
	sessions.Session
}

func Default(ctx *gin.Context) NoobSession {
	return NoobSession{sessions.Default(ctx)}
}

func (s NoobSession) SetM(data map[string]interface{}) {
	for k, v := range data {
		s.Set(k, v)
	}
}

func (s NoobSession) IsLoggedIn() bool {
	_, ok := s.Get("username").(string)
	return ok
}

func (s NoobSession) IsAdmin() bool {
	role, ok := s.Get("role").(string)
	return ok && role == "admin"
}

func (s NoobSession) Username() string {
	usr, ok := s.Get("username").(string)
	if !ok {
		return ""
	}

	return usr
}

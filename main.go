package main

import (
	"htmx/router"
)

func main() {
	//fixture.RunFixtures("fixtures", database.DatabaseParams{Name: "htmx"})

	// opt, err := redis.ParseURL("redis://redis:6379")
	// if err != nil {
	// 	log.Println("error connection redis parse")
	// }
	// redisClient := redis.NewClient(opt)
	// if err := redisClient.Ping(); err.String() != "ping: PONG" {
	// 	log.Println("error redis connection init: " + err.String())
	// } else {
	// 	log.Println("redis connection worked")
	// }

	router.InitRouter()
}

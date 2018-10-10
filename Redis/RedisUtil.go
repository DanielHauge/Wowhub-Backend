package Redis

import (
	"github.com/go-redis/redis"
	"log"
)



var Addr string = "localhost"
var Port string = ":6379"
var Password string = ""
var DB int = 0


func CanIConnect() error{
	client := redis.NewClient(&redis.Options{
		Addr: Addr+Port,
		Password: Password,
		DB: DB,
	})
	_, e := client.Ping().Result()

	return e
}

func DoesKeyExist(key string) bool{
	client := redis.NewClient(&redis.Options{
		Addr: Addr+Port,
		Password: Password,
		DB: DB,
	})
	d, e := client.Exists(key).Result()
	if e != nil{
		log.Println(e.Error())
	}

	if d == 1{
		return true
	} else {
		return false
	}
}

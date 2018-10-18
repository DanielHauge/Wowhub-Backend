package Redis

import (
	log "../Logrus"
	"github.com/go-redis/redis"
	"os"
)

var Addr = os.Getenv("CONNECTION_STRING")
var Port = ":6379"
var Password = ""
var DB = 0

// TODO: If availability ever becomes a problem, look into ClusterClient.
// TODO: If redis becomes cache only and availability becomes a problem, look into Ring for multiple redis servers.

func CanIConnect() error {
	client := redis.NewClient(&redis.Options{
		Addr:     Addr + Port,
		Password: Password,
		DB:       DB,
	})
	_, e := client.Ping().Result()

	return e
}

func DoesKeyExist(key string) bool {
	client := redis.NewClient(&redis.Options{
		Addr:     Addr + Port,
		Password: Password,
		DB:       DB,
	})
	d, e := client.Exists(key).Result()
	if e != nil {
		log.WithLocation().WithError(e).Error("Hov!")
	}

	if d == 1 {
		return true
	} else {
		return false
	}
}

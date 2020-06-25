package redis

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567089")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}


func Populate()  {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	categories := []string{"music","sports","news","food"}
	ids := make([]string, 4)
	log.Println("Adding ids : ")
	for i,v := range categories {
		ids[i] = RandStringRunes(7)
		log.Printf("%v\n", ids[i])
		_, err = conn.Do("HMSET", "video:"+ids[i], "title", "title:"+ RandStringRunes(7), "category" , v , "likes", 0)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Finished Populating Redis with 4 entries !!")
}

package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567089")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}


func Populate(numRecords int)  {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error while creating connection with redis")
		panic(err)
	}
	defer conn.Close()
	categories := []string{"music","sports","news","food"}
	ids := make([]string, numRecords)
	fmt.Println("Adding ids : ")
	for i,v := range categories {
		ids[i] = RandStringRunes(7)
		fmt.Printf("video:%v\n", ids[i])
		_, err = conn.Do("HSET", "video:"+ids[i], "title", RandStringRunes(7), "category" , v , "likes", 0)
		if err != nil {
			fmt.Println("Error while populating Redis")
			panic(err)
		}
	}
	fmt.Printf("Finished Populating Redis with %v entries\n", numRecords)
}

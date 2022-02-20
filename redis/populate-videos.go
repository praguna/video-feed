package redis

import (
	"fmt"
	"math/rand"

	"github.com/gomodule/redigo/redis"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567089")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetRandomCategory() int {
	categories := []string{"music", "sports", "news", "food"}
	randomIndex := rand.Intn(len(categories))
	fmt.Println(randomIndex)
	return 5
}

func Populate(numRecords int) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error while creating connection with redis")
		panic(err)
	}
	defer conn.Close()
	//
	ids := make([]string, numRecords)
	fmt.Println("Adding ids : ")
	// range categories is only four, so that it why we only get 4 records,
	for i := 0; i < numRecords; i++ {
		ids[i] = RandStringRunes(7)
		fmt.Printf("video:%v\n", ids[i])
		_, err = conn.Do("HSET", "video:"+ids[i], "title", RandStringRunes(7), GetRandomCategory(), "test", "likes", rand.Intn(50))
		if err != nil {
			fmt.Println("Error while populating Redis")
			panic(err)
		}
	}
	fmt.Printf("Finished Populating Redis with %v entries\n", numRecords)
}

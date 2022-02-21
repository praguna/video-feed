package redis

type Message struct {
	message string `redis:"messages"`
}

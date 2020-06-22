package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

var tmp []string

func main() {
	csv := os.Args[1]
	redis := os.Args[2]
	fmt.Println(csv, redis)
	redisconn := RedisInit(redis)
	defer (*redisconn).Close()

	channel := make(chan []string)

	start := time.Now()

	// ReadCSVtoRedis(csv, redisconn)

	go ReadCSVtoRedis(csv, channel)

	for v := range channel {
		tmp = v
		// fmt.Println(v)
	}

	elapsed := time.Now().Sub(start)
	fmt.Println("it takes ", elapsed)
}

func RedisInit(conn string) *redis.Conn {
	c, err := redis.Dial("tcp", conn)
	if err != nil {
		fmt.Println("Connect to redis error", err)
	}

	return &c
}

func ReadCSVtoRedis(fileName string, list chan<- []string) {

	fs, err := os.Open(fileName)
	defer fs.Close()
	if err != nil {
		fmt.Println("read csv error")
	}
	r := csv.NewReader(fs)
	index := 0
	for {
		index++
		row, err := r.Read()
		if err == io.EOF {
			close(list)
			break
		}
		if err != nil && err != io.EOF {
			fmt.Println("reading error")
		}
		// _, err = (*redisconn).Do("SET", "mykey", row)
		list <- row
		if err != nil {
			fmt.Println("redis set failed:", err)
		}
	}
	fmt.Println("Done for ", index)
}

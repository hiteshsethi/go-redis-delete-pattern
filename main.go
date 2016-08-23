package main

import "os"
import (
	"fmt"
)

func main() {
	args := os.Args

	if (len(args) != 2 || args[1] == "") {
		fmt.Println("Improper Args")
		return
	}

	pattern := args[1]

	redis := RedisStore{}

	pattern = "*" + pattern + "*"

	values, err := redis.Strings(redis.Run("keys", pattern))

	if err != nil {
		panic(err)
	}

	totalDeleted := 0

	for _, val := range values {
		_, err = redis.Run("DEL", val)

		if (err == nil) {
			totalDeleted += 1
		}
	}

	fmt.Println("Total Deleted Keys :", totalDeleted)

}
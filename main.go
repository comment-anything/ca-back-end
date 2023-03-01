package main

import (
	"fmt"

	"github.com/comment-anything/ca-back-end/config"
	"github.com/comment-anything/ca-back-end/server"
)

func main() {
	fmt.Println("Start")
	err := config.Vals.Load(".env")
	if err != nil {
		fmt.Printf("There was an error parsing environment variables: %s", err.Error())
	} else {
		s, err := server.New()
		if err == nil {
			s.Start(true)
		} else {
			fmt.Printf("\nError starting server: %s", err.Error())
		}
	}
}

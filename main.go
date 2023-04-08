package main

import (
	"flag"
	"fmt"

	"github.com/comment-anything/ca-back-end/config"
	"github.com/comment-anything/ca-back-end/server"
)

func main() {

	noCliflag := flag.Bool("nocli", false, "Starts server without a cli. Use for non-headless deployments where more information is desired.")

	customEnvVars := flag.String("env", ".env", "Sets a custom path for the server environment variables.")

	dockerMode := flag.Bool("docker", false, "Starts server in docker mode. In this mode, a different connection string is generated based on .env variables for when the server is running in a docker container.")

	flag.Parse()

	fmt.Println("Start")
	err := config.Vals.Load(*customEnvVars, *dockerMode)
	if err != nil {
		fmt.Printf("There was an error parsing environment variables: %s", err.Error())
	} else {
		s, err := server.New()
		if err == nil {
			s.Start(!*noCliflag)
		} else {
			fmt.Printf("\nError starting server: %s", err.Error())
		}
	}
	fmt.Printf("\nServer finished.")
}

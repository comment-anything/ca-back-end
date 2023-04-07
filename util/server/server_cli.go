package server

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/comment-anything/ca-back-end/config"
)

var cli_is_running bool = true

var serv *Server

// CLIBegin starts the server command-line-interface loop. This runs in a different thread than the one which is actually acting as a server, listening to user responses. The CLI provides various commands to view information about the server and change settings on the server.
func (s *Server) CLIBegin() {
	serv = s
	// If user presses a key, we terminate the server.
	reader := bufio.NewReader(os.Stdin)
	for cli_is_running { // loop until user runs exit command
		fmt.Printf(getStatusHeader())
		userinput, _ := reader.ReadString('\n')
		userinput = userinput[:len(userinput)-1] /* chop the newline and EOF chars */
		parseCommand(userinput)
	}
}

func parseCommand(inp string) {
	if inp == "exit" || inp == "stop" {
		fmt.Println("\nExiting server.")
		err := serv.httpServer.Shutdown(context.Background())
		fmt.Printf("Server shutdown gracefully with error: %v\n", err)
		cli_is_running = false
	} else if inp == "log" {
		config.Vals.Server.DoesLogAll = !config.Vals.Server.DoesLogAll
		fmt.Printf("\nToggled Logging to: %v\n\n", config.Vals.Server.DoesLogAll)
	} else if inp == "user count" {
		fmt.Printf("\n%s\n\n", serv.users.GetUserCountString())
	} else if inp == "user list" {
		fmt.Printf("\n%s\n\n", serv.users.GetUserListString())
	} else if inp == "pages count" {
		fmt.Printf("\n%s\n\n", serv.PageManager.GetPageManagerCountString())
	} else if inp == "pages list" {
		fmt.Printf("\n%s\n\n", serv.PageManager.GetPagesListString())
	} else if len(inp) >= 9 && inp[0:9] == "page info" {
		fmt.Printf("\n%s\n\n", serv.PageManager.GetPageInfo(inp))
	} else {
		fmt.Printf(getHelp())

	}
}

// this will print at the start of a line.
func getStatusHeader() string {
	return fmt.Sprintf("[%v]>> ", serv.httpServer.Addr)
}

func getHelp() string {
	return fmt.Sprintf(
		"\n" + ` Available commands
   -- cmd --      --- effect ---
   stop/exit      stop the server and exit the cli.
   log            toggle logging
   user count     get a count of the users
   user list      get a list of the logged-in users
   pages count    get a count of the pages
   pages list	  list the loaded pages
   page info n    gets page info for page with id n
` + "\n")
}

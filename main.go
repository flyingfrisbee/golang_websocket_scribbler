package main

import (
	"GithubRepository/golang_websocket_scribbler/cmd"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd.Run()
}

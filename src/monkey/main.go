package main

import (
	"fmt"
	"os"
	"os/user"

	"go-compiler/src/monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nHello %s! This is the Monkey Programming Language!\n", user.Username)
	fmt.Println("Statements need a semicolon to end; enter `exit()` or CTRL-d (i.e. EOF) to exit. Synatx: https://monkeylang.org")
	fmt.Printf("\n")

	repl.REPL(os.Stdin, os.Stdout)
}

package main

import (
	"avm/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This Abstract VM\n", usr.Username)
	repl.Start(os.Stdin, os.Stdout)

}

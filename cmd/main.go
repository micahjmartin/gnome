package main

import (
	"fmt"
	"os"

	"github.com/nullmonk/gnome"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("USAGE: main <script.eld>")
		os.Exit(1)
	}
	gnome.Run()
}

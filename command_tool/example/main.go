package main

import (
	"io"
	"os"

	"github.com/leilei3167/ci/command_tool/cmd"
)

func main() {
	c := cmd.NewXXXComand(os.Stdout, io.Discard, io.Discard)
	c.Execute()
}

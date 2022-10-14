package main

import (
	"io"
	"os"

	"github.com/leilei3167/ci/command_tool/cmd"
)

func main() {
	c := cmd.NewXXXComand(os.Stdout, io.Discard, io.Discard)
	c.Execute()

	// 左对齐,并且至少要有50的宽度,如果没有-,将不会左对齐
	// tmp := fmt.Sprintf("%%-%ds", 50)
	// s := fmt.Sprintf(tmp, testString) + "tag\n"
	// fmt.Print(s)

	// //%.Xs 限制最大宽度为2
	// tmp1 := fmt.Sprintf("%%.%ds\n", 2)
	// fmt.Printf(tmp1, testString)

	// //%x.xs 最小为5,最大为7
	// tmp2 := fmt.Sprintf("%%%d.%ds\n", 5, 7)
	// fmt.Printf(tmp2, testString)

	// //%x.xs
	// tmp3 := fmt.Sprintf("%%-%d.%ds\n", 5, 3) //%-5.3s 大于3则截断,保留5宽度
	// fmt.Printf(tmp3, testString)
}

const (
	testString = "THIS IS A TEST COMMAND"
)

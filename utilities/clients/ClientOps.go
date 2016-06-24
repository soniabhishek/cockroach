package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	userName := flag.String("u", "foo", "a string")
	password := flag.String("u", "foo", "a string")
	projectName := flag.String("u", "foo", "a string")
	projectDesc := flag.String("u", "foo", "a string")
	projectDesc := flag.String("u", "foo", "a string")
	url := flag.Int("numb", 42, "an int")
	boolPtr := flag.Bool("fork", false, "a bool")

	flag.Parse()

	getArgs()
}

func getArgs() {

	fmt.Println(len(os.Args), os.Args)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}

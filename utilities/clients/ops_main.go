package main

import (
	"fmt"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/utilities/clients/models"
	"gitlab.com/playment-main/angel/utilities/clients/operations"
)

/*-func main() {
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
}*/

func main() {
	fmt.Println([]byte("password"))
	/*
		Using development environment
		[112 97 115 115 119 111 114 100]
	*/
}
func some() {
	cl := utilModels.Client{
		UserName:       "experiment",
		Password:       []byte("exprpassword"),
		ClientId:       uuid.Nil,
		ClientSecretId: uuid.Nil,
		ProjectId:      uuid.Nil,
		ProjectLabel:   "exprLabel",
		ProjectName:    "exprName",
		Url:            "https://www.google.com",
		Header:         models.JsonFake{"1": "One"},
		Steps:          nil,

		Gender:    "m",
		FirstName: "firstExpr",
		LastName:  "secondExpr",
		Phone:     "7788996655",
	}

	service := operations.Service{}
	status, err := service.CreateClient(cl)
	fmt.Println(status, err)
}

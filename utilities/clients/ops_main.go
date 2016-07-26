package main

import (
	"flag"

	"gitlab.com/playment-main/angel/app/api/auther"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/utilities/clients/models"
	"gitlab.com/playment-main/angel/utilities/clients/operations"
	"gitlab.com/playment-main/angel/utilities/clients/validator"
)

func main() {
	/* Important */
	userName := flag.String("username", "", "username")
	password := flag.String("password", "", "password")
	projectLabel := flag.String("projectLabel", "", "provide project label")
	projectName := flag.String("projectName", "", "provide project name")
	url := flag.String("url", "", "url")
	headerStr := flag.String("header", "", "a json")
	header := models.JsonFake{}
	header.Scan(*headerStr)

	/* Optional */
	gender := flag.String("gender", "", "Gender [optional]")
	firstName := flag.String("firstname", "", "First Name [optional]")
	lastName := flag.String("lastname", "", "Second Name [optional]")
	phone := flag.String("phone", "", "Phone Number [optional]")
	flag.Parse()

	obj := utilModels.Client{
		UserName:       *userName,
		Password:       *password,
		ClientId:       uuid.Nil,
		ClientSecretId: uuid.Nil,
		ProjectId:      uuid.Nil,
		ProjectLabel:   *projectLabel,
		ProjectName:    *projectName,
		Url:            *url,
		Header:         header,
		Steps:          nil,

		Gender:    *gender,
		FirstName: *firstName,
		LastName:  *lastName,
		Phone:     *phone,
	}

	err := validator.ValidateInput(obj)
	if err != nil {
		validator.ShowErrorResponse(err)
		return
	}

	service := operations.Service{}

	_, err = service.CreateClient(&obj)
	if err == nil {
		result := models.JsonFake{
			"success": true,
			"userdetails": utilModels.Client{
				UserName:        obj.UserName,
				ClientId:        obj.ClientId,
				ClientSecretStr: auther.StdProdAuther.GetAPIKey(obj.ClientSecretId),
				ProjectId:       obj.ProjectId,
				ProjectLabel:    obj.ProjectLabel,
				ProjectName:     obj.ProjectName,
				Url:             obj.Url,
				Header:          obj.Header,
				Steps:           obj.Steps,

				Gender:    obj.Gender,
				FirstName: obj.FirstName,
				LastName:  obj.LastName,
				Phone:     obj.Phone,
			},
		}
		plog.Info(result.StringPretty())
	} else {
		plog.Error("Error while creating user: ", err)
		validator.ShowErrorResponse(err)
	}
}

package main

import (
	"brcaidsurvey/pkg/model"
	"fmt"
	"os"
)

func main() {
	fmt.Println("In DB Prep")
	err := model.InitDB()
	if err != nil {
		fmt.Println("Bad things " + err.Error())
	}
	var path string
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		dir, _ := os.Getwd()
		path = fmt.Sprintf("%s/formdata/survey.json", dir)
	}

	err = model.LoadFormDataIntoDB(path)
	if err != nil {
		fmt.Println("Bad things " + err.Error())
	}

	adminUser := model.User{UserUUID: "0", UserID: "admin", PasswordHash: model.HashPassword("admin"), Name: "admin"}
	err = model.PutUser(&adminUser)
	if err != nil {
		fmt.Println("Trouble adding admin user")
	}

}

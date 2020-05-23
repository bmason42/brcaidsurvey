package main

import (
	"brcaidsurvey/pkg/db"
	"fmt"
	"os"
)

func main(){
	fmt.Println("In DB Prep")
	err:= db.InitDB()
	if err != nil{
		fmt.Println("Bad things " + err.Error())
	}
	var path string
	if len(os.Args) >1{
		path=os.Args[1]
	}else{
		dir, _ := os.Getwd()
		path = fmt.Sprintf("%s/formdata/survey.json", dir)
	}

	err=db.LoadFormDataIntoDB(path)
	if err != nil{
		fmt.Println("Bad things " + err.Error())
	}

}

package main

import (
	"brcaidsurvey/pkg"
	"fmt"
)

func main(){
	fmt.Println("Hello")
	err:= pkg.InitDB()
	if err != nil{
		fmt.Println("Bad things " + err.Error())
	}

}

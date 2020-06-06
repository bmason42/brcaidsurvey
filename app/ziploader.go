/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package main

import (
	"brcaidsurvey/pkg/db"
	"brcaidsurvey/pkg/model"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Need a zip code csv")
		os.Exit(1)
	}
	db.InitDB()
	connection, err := db.GetDBConnection()
	if err != nil {
		fmt.Printf("DB Open Error %s", err.Error())
		os.Exit(1)
	}
	in, e := os.Open(os.Args[1])
	if e != nil {
		fmt.Printf("Error  reading file %s\n", os.Args[1])
		os.Exit(2)
	}
	zipMap := make(map[string]string)
	reader := csv.NewReader(in)
	for {
		rec, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("Error  reading file %s\n", os.Args[1])
				os.Exit(2)
			}
		}
		//skip deommentioned numbers
		if rec[15] == "TRUE" {
			continue
		}
		if rec[1] == "Zipcode" {
			//skip title row
			continue
		}
		var zip model.ZipCode
		zip.ZipCode = rec[1]
		zip.City = rec[3]
		zip.State = rec[4]
		zip.Lat, _ = strconv.ParseFloat(rec[6], 32)
		zip.Long, _ = strconv.ParseFloat(rec[7], 32)
		zip.Xaxis, _ = strconv.ParseFloat(rec[8], 32)
		zip.Yaxis, _ = strconv.ParseFloat(rec[9], 32)
		zip.Zaxis, _ = strconv.ParseFloat(rec[10], 32)

		_, ok := zipMap[zip.ZipCode]
		if !ok {
			zipMap[zip.ZipCode] = ""
			connection.Create(&zip)
		}

	}
}

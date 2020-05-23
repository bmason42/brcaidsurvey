package db

import (
	"brcaidsurvey/pkg/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func LoadFormDataIntoDB(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	var data model.FormData
	file, e := os.Open(path)
	if e != nil {
		return e
	}
	var bits []byte
	bits, e = ioutil.ReadAll(file)
	if e != nil {
		return e
	}
	e = json.Unmarshal(bits, &data)

	db, e := GetDBConnection()
	if e != nil {
		return e
	}
	delStm := fmt.Sprintf("delete from %s;", db.NewScope(&model.RegionInfo{}).TableName())
	db.Exec(delStm)
	delStm = fmt.Sprintf("delete from %s;", db.NewScope(&model.SupportConcern{}).TableName())
	db.Exec(delStm)

	for _, x := range data.BmRegions {
		db.Create(x)
	}
	for _, x := range data.Concerns {
		db.Create(x)
	}

	return e
}

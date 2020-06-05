package db

import (
	"brcaidsurvey/pkg/model"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func LoadFormDataIntoDB(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	var data model.FormData
	data.Concerns = make([]model.SupportConcern, 0)
	data.BmRegions = make([]model.RegionInfo, 0)
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
		if db.Error != nil {
			log.Infof("Error with record %s", x)
		}

	}
	for _, x := range data.Concerns {
		db.Create(x)
		if db.Error != nil {
			log.Infof("Error with record %s", x)
		}
	}

	return e
}

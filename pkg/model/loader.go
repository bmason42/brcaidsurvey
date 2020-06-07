/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package model

import (
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
	var data FormData
	data.Skills = make([]Skill, 0)
	data.BmRegions = make([]RegionInfo, 0)
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
	delStm := fmt.Sprintf("delete from %s;", db.NewScope(&RegionInfo{}).TableName())
	db.Exec(delStm)
	delStm = fmt.Sprintf("delete from %s;", db.NewScope(&Skill{}).TableName())
	db.Exec(delStm)
	for _, x := range data.BmRegions {
		db.Create(x)
		if db.Error != nil {
			log.Infof("Error with record %s", x)
		}

	}
	for _, x := range data.Skills {
		db.Create(x)
		if db.Error != nil {
			log.Infof("Error with record %s", x)
		}
	}

	return e
}

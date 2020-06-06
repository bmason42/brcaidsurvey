/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package main

import (
	"brcaidsurvey/pkg/apiimpl"
	"brcaidsurvey/pkg/model"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := model.InitModel()
	if err != nil {
		log.Error("No config found " + err.Error())
		panic("Cannot continue")
	}
	apiimpl.RunServer()
}

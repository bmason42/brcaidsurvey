/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package main

import (
	"brcaidsurvey/pkg/apiimpl"
	"brcaidsurvey/pkg/model"
	"testing"
)

func Test_main(t *testing.T) {
	err := model.InitModel()
	if err != nil {
		panic(err)
	}
	apiimpl.RunServer()
}

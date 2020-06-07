/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package tests

import (
	"brcaidsurvey/pkg/apiimpl"
	"brcaidsurvey/pkg/generated/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContactConversions(t *testing.T) {
	var api v1.SurveyContact
	api.NeedHelpNow = true
	api.OfferedSkills = []string{"skill1", "skill2"}
	api.ConactInfo = v1.SurveyContactPii{Name: "bob", Email: "bob@gmail.com", Phone: "503-555-1212", Zip: "62650", PreferedContact: "email"}

	model := apiimpl.SurveyContactApiToModel(&api)
	api2, err := apiimpl.SurveyContactModelToApi(model)
	assert.Nil(t, err)
	assert.NotNil(t, api2)
	assert.Equal(t, api.ConactInfo, api2.ConactInfo)
}

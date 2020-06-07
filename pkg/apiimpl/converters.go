/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package apiimpl

import (
	"brcaidsurvey/pkg/generated/v1"
	"brcaidsurvey/pkg/model"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
)

func SurveyContactApiToModel(contact *v1.SurveyContact) *model.SurveyContact {
	var ret model.SurveyContact
	ret.SurveyContactID = uuid.New().String()
	ret.EncryptionVersion = model.ENCRYPT_VERSION_ONE
	ret.NeedHelpNow = contact.NeedHelpNow
	ret.OfferingHelp = contact.OfferingHelp
	ret.RequestingHelp = contact.RequestingHelp
	ret.ZipCode = contact.Zip
	first := true
	for _, x := range contact.NeededSkills {
		if first {
			first = false
		} else {
			ret.RequestedSkills = ret.RequestedSkills + "|"
		}
		ret.RequestedSkills = ret.RequestedSkills + x
	}
	first = true
	for _, x := range contact.OfferedSkills {
		if first {
			first = false
		} else {
			ret.OfferedSkills = ret.OfferedSkills + "|"
		}
		ret.OfferedSkills = ret.OfferedSkills + x

	}

	cipher := model.PlainStructToCipher(ret.SurveyContactID, &contact.ConactInfo)
	bits, _ := json.Marshal(&cipher)
	ret.PII = string(bits)
	return &ret
}
func SurveyContactModelToApi(contact *model.SurveyContact, includePII bool) (*v1.SurveyContact, error) {
	var ret v1.SurveyContact
	var cipherRecord model.CipherRecord
	if includePII {
		err := json.Unmarshal([]byte(contact.PII), &cipherRecord)
		if err != nil {
			return nil, err
		}
		err = model.CipherRecordToPlainRecord(&cipherRecord, &ret.ConactInfo)
		if err != nil {
			return nil, err
		}
	}

	ret.OfferedSkills = strings.Split(contact.OfferedSkills, "|")
	ret.NeededSkills = strings.Split(contact.RequestedSkills, "|")
	ret.RequestingHelp = contact.RequestingHelp
	ret.OfferingHelp = contact.OfferingHelp
	ret.NeedHelpNow = contact.NeedHelpNow
	ret.Zip = contact.ZipCode

	return &ret, nil
}

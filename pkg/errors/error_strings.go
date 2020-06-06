/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package errors

import "fmt"

var errorStrings map[string]map[string]string

func init() {
	errorStrings = make(map[string]map[string]string)
	m := mkUSMap()
	errorStrings["en"] = m
}

func GetErrorString(locale, errcode string) string {
	if _, ok := errorStrings[locale]; !ok {
		locale = "en"
	}
	if _, ok := errorStrings[locale][errcode]; !ok {
		errcode = ERROR_CODE_UNKNOWN
	}
	return errorStrings[locale][errcode]
}
func mkUSMap() map[string]string {
	ret := make(map[string]string)
	ret[fmt.Sprintf("%s.%s", NETWOR_ERROR, NETWORK_ERROR_INVALID_URL)] = "Invalid URL provided"
	ret[fmt.Sprintf("%s.%s", NETWOR_ERROR, NETWORK_ERROR_CONNECT_FAIL)] = "Failed to connect to external service"
	ret[fmt.Sprintf("%s.%s", OCERROR_ERROR, ERROR_CODE_UNKNOWN)] = "An unknown error occurred. "
	ret[fmt.Sprintf("%s.%s", OCERROR_ERROR, ERROR_CODE_INVALID_USER_INPUT)] = "Invalid user data. "
	ret[fmt.Sprintf("%s.%s", OCERROR_ERROR, ERROR_CODE_SYSTEM_NOT_INITIALIZED)] = "The model of the system has not been initialized. "
	ret[fmt.Sprintf("%s.%s", OCERROR_ERROR, ERROR_CODE_NOT_VALID_NUMBER)] = "The provided phone number does not appear to be valid.  Please enter a valid number. "

	return ret
}

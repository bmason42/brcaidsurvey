/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package errors

import "fmt"

//opencrisis2 errors
const (
	ERROR_CODE_UNKNOWN                = "unknown"
	ERROR_CODE_INVALID_USER_INPUT     = "invalidrequest"
	ERROR_CODE_SYSTEM_NOT_INITIALIZED = "noinit"
	ERROR_CODE_NOT_VALID_NUMBER       = "numbernotvalid"
)

//network errors
const (
	NETWORK_ERROR_INVALID_URL  = "invalidurl"
	NETWORK_ERROR_CONNECT_FAIL = "connectfailed"
)

const (
	OCERROR_ERROR = "opencrisis2"
	NETWOR_ERROR  = "network"
)

type InternalError struct {
	Subsystem      string
	SubSystemError string
	Params         map[string]string
}

func (t *InternalError) Error() string {
	return fmt.Sprintf("%s.%s", t.Subsystem, t.SubSystemError)
}
func (t *InternalError) ErrorCode() string {
	return fmt.Sprintf("%s.%s", t.Subsystem, t.SubSystemError)
}
func NewInernalError(subSystem, code string, params map[string]string) *InternalError {
	var x InternalError
	x.Subsystem = subSystem
	x.SubSystemError = code
	x.Params = params
	return &x

}

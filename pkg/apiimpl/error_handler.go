/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package apiimpl

import (
	"context"
	"fmt"
	"github.com/bmason42/opencrisisline2/pkg/errors"
	"github.com/bmason42/opencrisisline2/pkg/generated/v1"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func MkErrorResponse(subsystem, code string, c context.Context, parameters map[string]string) *v1.ErrorResponse {
	locale := getLocale(c)
	message := errors.GetErrorString(locale, fmt.Sprintf("%s.%s", subsystem, code))
	return newErrorResponse(subsystem, code, message, parameters)
}
func newErrorResponse(subsystem, code string, message string, parameters map[string]string) *v1.ErrorResponse {
	if parameters == nil {
		parameters = make(map[string]string)
	}
	fullCode := fmt.Sprintf("%s.%s", subsystem, code)
	return &v1.ErrorResponse{
		Code:       fullCode,
		Message:    message,
		Parameters: parameters,
	}
}

func handleErrors(ctx context.Context, errs ...error) (httpStatusCode int, resps []*v1.ErrorResponse) {
	for _, e := range errs {
		c, resp := handleError(ctx, e)
		if c > httpStatusCode {
			httpStatusCode = c
		}
		resps = append(resps, resp)
	}
	return
}
func getLocale(ctx context.Context) string {
	return "en"
}
func handleError(ctx context.Context, e error) (httpStatusCode int, resp *v1.ErrorResponse) {
	locale := getLocale(ctx)
	// Defaulting since there is currently no other supported languages

	switch e.(type) {
	case *errors.InternalError:
		x := e.(*errors.InternalError)
		log.Errorf("internal error %s", e.Error())
		httpStatusCode = http.StatusBadRequest

		resp = newErrorResponse(x.Subsystem, x.SubSystemError, errors.GetErrorString(locale, x.ErrorCode()), x.Params)

	default:
		log.Errorf("An unhandled error occurred: %T: %s", e, e.Error())
		httpStatusCode = http.StatusInternalServerError
		params := make(map[string]string)
		params["detail"] = e.Error()
		resp = newErrorResponse(errors.OCERROR_ERROR, errors.ERROR_CODE_UNKNOWN, errors.GetErrorString(locale, errors.ERROR_CODE_UNKNOWN), params)
	}
	return
}

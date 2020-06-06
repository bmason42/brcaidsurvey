/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package apiimpl

import (
	"fmt"
	"github.com/bmason42/opencrisisline2/pkg/errors"
	"github.com/bmason42/opencrisisline2/pkg/generated/v1"
	"github.com/bmason42/opencrisisline2/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"strings"
	"unicode"
)

func sendError(c *gin.Context, err error) {
	httpCode, errs := handleErrors(c, err)
	c.JSON(httpCode, errs)
}

func aboutGetUnversioned(c *gin.Context) {
	var resp v1.AboutResponse
	resp.AppVersion = "1.0"
	resp.ApiVersions = make([]string, 0)
	resp.ApiVersions = append(resp.ApiVersions, "v1")

	c.JSON(http.StatusOK, resp)
	log.Println("In about")
}

type Message struct {
	Body string
}
type Response struct {
	Message Message
}

func callbackHandler(c *gin.Context) {
	c.Request.ParseForm()
	for key, value := range c.Request.Form {
		fmt.Printf("%s = %s \n", key, value)
	}
	pinToFind := c.Request.Form.Get("Body")
	layer, e := model.GetPersistenceLayer()
	if e != nil {
		log.Errorf("Unable to access persistence layer in Callback handler %s \n", e.Error())
		c.XML(500, "")
		return
	}
	fromNumber := c.Request.Form.Get("From")
	if len(fromNumber) < 10 {
		log.Error("Unexpected  From Number format %s \n", fromNumber)
		c.XML(500, "")
		return
	}

	//get rid of country code and any pluses
	fromNumber = fromNumber[len(fromNumber)-10:]
	record, err := layer.FetchSupportRequestByPhone(fromNumber)
	if err != nil {
		log.Errorf("Error looking up records in callback %s \n", err.Error())
		c.XML(500, "")
		return
	}

	if record == nil {
		log.Warnf("Unable to find record for phone number %s \n", fromNumber)
		c.XML(500, "")
		return
	}

	var r Response
	pinToFind = strings.TrimSpace(pinToFind)
	if record.AuthPin != pinToFind {
		r.Message.Body = "PIN number does not match"
	} else {
		r.Message.Body = "Expect a text from someone shortly"
		ms := model.NewMessageSystem()
		ms.AsyncSendText(model.GetConfig().DutyNumber, fmt.Sprintf("%s, %s, %s", record.Data.PhoneNumber, record.Data.CallerName, record.Data.Message))
	}

	c.XML(200, r)
}

func postHandler(c *gin.Context) {
	var help v1.HelpRequest
	err := c.ShouldBindJSON(&help)
	if err != nil {
		weberr := MkErrorResponse(errors.OCERROR_ERROR, errors.ERROR_CODE_INVALID_USER_INPUT, c, map[string]string{"data": err.Error()})
		c.JSON(400, weberr)
	}
	scrubbedPhone := make([]rune, 0)
	for _, r := range []rune(help.PhoneNumber) {
		if unicode.IsNumber(r) {
			scrubbedPhone = append(scrubbedPhone, r)
		}
	}
	phoneString := string(scrubbedPhone)
	if len(phoneString) < 10 {
		weberr := MkErrorResponse(errors.OCERROR_ERROR, errors.ERROR_CODE_NOT_VALID_NUMBER, c, map[string]string{"data": help.PhoneNumber})
		c.JSON(400, weberr)
	}
	help.PhoneNumber = phoneString[len(phoneString)-10:]

	var supportReq model.SupportRequest
	supportReq.RequestID = uuid.New().String()
	pinNumber := rand.Uint32() % 10000
	supportReq.AuthPin = fmt.Sprintf("%04d", pinNumber)
	supportReq.Data = help
	layer, _ := model.GetPersistenceLayer()
	layer.SaveSupportRequest(&supportReq)

	ms := model.NewMessageSystem()
	go ms.AsyncSendText(help.PhoneNumber, "Please reply to this text with ONLY the PIN Number from the Web Site")
	var resp v1.HelpResponse
	resp.AuthPin = supportReq.AuthPin
	resp.RequestID = supportReq.RequestID

	c.JSON(201, &resp)

}
func healthCheckGetUnversioned(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
func swaggerUIGetHandler(c *gin.Context) {
	c.Redirect(302, "/opencrisisline2/swaggerui/index_v1.html")
}

/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package apiimpl

import (
	"brcaidsurvey/pkg/generated/v1"
	"brcaidsurvey/pkg/model"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
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

func healthCheckGetUnversioned(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
func swaggerUIGetHandler(c *gin.Context) {
	c.Redirect(302, "/brcaid/swaggerui/index_v1.html")
}
func handleSurveyGet(c *gin.Context) {
	modelData, err := model.FetchSurveyData()
	if err != nil {
		c.JSON(handleError(c, err))
		return
	}
	ret := make([]v1.SurveyContact, len(modelData))
	for i, x := range modelData {
		tmp, _ := SurveyContactModelToApi(&x, false)
		ret[i] = *tmp
	}
	c.JSON(200, &ret)
}
func handleSurveyPost(c *gin.Context) {
	var req v1.SurveyContact
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(handleError(c, err))
		return
	}

	modelData := SurveyContactApiToModel(&req)
	err = model.PutSurveyContact(modelData)
	if err != nil {
		c.JSON(handleError(c, err))
		return
	}
	var ret v1.SurveyContactResp
	ret.RequestUUID = modelData.SurveyContactID

}
func loginHandler(c *gin.Context) {
	var data v1.Login
	err := c.ShouldBind(&data)
	if err != nil {
		c.JSON(handleError(c, err))
		return
	}
	b := model.ValidatePassword(data.UserID, data.Password)
	if b {
		session := model.MakeNewSession(data.UserID)
		resp := v1.LoginResponse{AutoToken: session.SessionID}
		c.JSON(201, &resp)
	} else {
		c.JSON(401, "")
	}
}
func logoutHandler(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")
	model.RemoveSession(authToken)
	c.JSON(204, "")
}
func userGetHandler(c *gin.Context) {
	users, e := model.FetchUsers()
	if e != nil {
		c.JSON(handleError(c, e))
		return
	}
	ret := make([]v1.User, len(users))
	for i := range users {
		ret[i].UserID = users[i].UserID
		ret[i].UserUUID = users[i].UserUUID
		ret[i].Phone = users[i].Phone
		ret[i].Email = users[i].Email
		//ret[i].=users[i].Name
	}
	c.JSON(200, ret)
}

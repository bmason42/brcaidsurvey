/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package tests

import (
	"brcaidsurvey/pkg/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestBasics(t *testing.T) {
	os.Setenv("unittests", "yes")
	model.InitModel()
	hash := model.HashPassword("bobo")
	u := model.User{UserUUID: uuid.New().String(), UserID: "admin", PasswordHash: hash, Name: "bob"}

	err := model.PutUser(&u)
	assert.Nil(t, err)

	userCopy, err := model.FetchUser(u.UserUUID)
	assert.Nil(t, err)
	assert.NotNil(t, userCopy)
	assert.Equal(t, u.UserUUID, userCopy.UserUUID)
	assert.Equal(t, u.UserID, userCopy.UserID)

	userCopy, err = model.FetchUserUserID(u.UserID)
	assert.Nil(t, err)
	assert.NotNil(t, userCopy)
	assert.Equal(t, u.UserUUID, userCopy.UserUUID)
	assert.Equal(t, u.UserID, userCopy.UserID)

	validPassowrd := model.ValidatePassword(u.UserID, "bobo")
	assert.Equal(t, true, validPassowrd)
	validPassowrd = model.ValidatePassword(u.UserID, "bobox")
	assert.Equal(t, false, validPassowrd)

	group := model.UserGroup{GroupUUID: uuid.New().String(), GroupName: "admins", GroupDescription: "those guys"}
	err = model.PutUserGroup(&group)
	assert.Nil(t, err)

	groupRanger := model.UserGroup{GroupUUID: uuid.New().String(), GroupName: "rangers", GroupDescription: "those rangers"}
	err = model.PutUserGroup(&groupRanger)
	assert.Nil(t, err)

	err = model.AddUserToGroup(u.UserUUID, group.GroupUUID)
	assert.Nil(t, err)
	err = model.AddUserToGroup(u.UserUUID, groupRanger.GroupUUID)
	assert.Nil(t, err)
	err = model.AddRoleToGroup(group.GroupUUID, model.ROLE_ID_ADMIN)
	assert.Nil(t, err)
	err = model.AddRoleToGroup(groupRanger.GroupUUID, model.ROLE_ID_RANGER)
	assert.Nil(t, err)
	err = model.AddRoleToGroup(groupRanger.GroupUUID, model.ROLE_ID_LICENSED_PROVIDER)
	assert.Nil(t, err)

	roles, err := model.FetchRoleForUser(u.UserUUID)
	assert.Nil(t, err)
	assert.NotNil(t, roles)
	assert.Equal(t, 3, len(roles))

}

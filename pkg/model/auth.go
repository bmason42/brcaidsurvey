/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package model

type RoleType string

const (
	ROLE_ID_ADMIN             RoleType = "auth.admin"
	ROLE_ID_RANGER            RoleType = "ranger"
	ROLE_ID_LICENSED_PROVIDER RoleType = "licensed.provider"
)

type User struct {
	UserUUID     string `gorm:"type: varchar(36);primary_key"`
	PasswordHash string
}

type UserGroup struct {
	GroupUUID        string `gorm:"type: varchar(36);primary_key"`
	GroupName        string
	GroupDescription string
}
type UserGroupX struct {
	UserUUID  string `gorm:"type: varchar(36);primary_key"`
	GroupUUID string `gorm:"type: varchar(36);primary_key"`
}

type Permission struct {
	GroupUUID string `gorm:"type: varchar(36);primary_key"`
	//roles are generated as needed to support user/group access
	RoleID string `gorm:"type: varchar(36);primary_key"`
}

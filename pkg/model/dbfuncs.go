/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package model

import (
	"fmt"
	_ "github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

var dbPath string

func InitDB() error {
	//left overs from mysql
	//home := os.Getenv("HOME")
	//dbPath = fmt.Sprintf("%s/brcaid-db/db.db", home)
	createDbORM()
	return nil
}

func GetDBConnection() (*gorm.DB, error) {
	//db, error := gorm.Open("sqlite3", dbPath)
	var db *gorm.DB
	var error error
	if os.Getenv("unittests") == "yes" {
		//for unit tests, we assume a local container
		db, error = gorm.Open("mysql", "root:howdy@tcp(localhost:3307)/brcaid?charset=utf8&parseTime=True&loc=Local")
	} else {
		dbuser := GetEnvVar("DB_USER", "root")
		dbPassword := GetEnvVar("DB_PASSWORD", "")
		dbHost := GetEnvVar("DB_HOST", "localhost:3306")
		connectArg := fmt.Sprintf("%s:%s@tcp(%s)/brcaid?charset=utf8&parseTime=True&loc=Local", dbuser, dbPassword, dbHost)
		db, error = gorm.Open("mysql", connectArg)
	}

	return db, error
}
func createDbORM() {
	db, _ := GetDBConnection()
	db.AutoMigrate(&RegionInfo{},
		&Skill{},
		&SurveyContact{},
		&SurveyResult{},
		&User{},
		&UserGroup{},
		&UserGroupX{},
		&Permission{},
		&ZipCode{})
	//db.Model(&model.SurveyResult{}).AddForeignKey("survey_contact_id", "survey_contacts(survey_contact_id)", "RESTRICT", "RESTRICT")

	//db.Raw("DROP INDEX event_idx on log_entry")
	//db.Raw("DROP INDEX tracking_idx on log_entry")
}

func PutUser(user *User) error {
	db, err := GetDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	var res *gorm.DB
	var existingUser User
	db.First(&existingUser, "user_uuid = ?", user.UserUUID)

	if existingUser.UserUUID == user.UserUUID {
		res = db.Save(user)
	} else {
		res = db.Create(user)
	}
	return res.Error
}
func PutUserGroup(group *UserGroup) error {
	db, err := GetDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	var res *gorm.DB
	var x UserGroup
	db.First(&x, "group_uuid = ?", group.GroupUUID)
	if x.GroupUUID == group.GroupUUID {
		res = db.Save(group)
	} else {
		res = db.Create(group)
	}
	return res.Error
}
func AddUserToGroup(userUUID, groupUUID string) error {
	x := UserGroupX{UserUUID: userUUID, GroupUUID: groupUUID}
	db, err := GetDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	var res *gorm.DB

	res = db.Create(&x)
	return res.Error
}
func AddRoleToGroup(groupUUID string, roleID RoleType) error {
	db, err := GetDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	var res *gorm.DB

	perm := Permission{GroupUUID: groupUUID, RoleID: roleID}
	res = db.Create(&perm)
	return res.Error
}

func FetchRoleForUser(userUUID string) (RoleMap, error) {
	ret := make(RoleMap, 0)
	db, err := GetDBConnection()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	var groups []UserGroupX
	res := db.Where("user_uuid=?", userUUID).Find(&groups)
	if res.Error != nil {
		return nil, res.Error
	}

	for _, groupX := range groups {
		var perms []Permission
		res = db.Where("group_uuid=?", groupX.GroupUUID).Find(&perms)
		if res.Error != nil {
			return nil, res.Error
		}
		for _, perm := range perms {
			ret[perm.RoleID] = true
		}
	}

	return ret, nil
}

func FetchUser(userUUID string) (*User, error) {
	db, err := GetDBConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var res *gorm.DB

	var user User
	res = db.Where("user_uuid=?", userUUID).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	var ret *User
	if user.UserUUID == userUUID {
		ret = &user
	}
	return ret, nil
}

func FetchUserUserID(userID string) (*User, error) {
	db, err := GetDBConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var res *gorm.DB

	var user User
	res = db.Where("user_id=?", userID).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	var ret *User
	if user.UserID == userID {
		ret = &user
	}
	return ret, nil
}

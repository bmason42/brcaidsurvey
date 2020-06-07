/*
 * Copyright (c) 2020.  This software is made for the Black Rock City Aid group and is provided AS IS with no support or liability under the Apache 2 license.
 */

package model

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type RoleType string
type RoleMap map[RoleType]bool

const (
	ROLE_ID_ADMIN             RoleType = "auth.admin"
	ROLE_ID_RANGER            RoleType = "ranger"
	ROLE_ID_LICENSED_PROVIDER RoleType = "licensed.provider"
)

type Session struct {
	UserID    string
	SessionID string
	LastTouch time.Time
	Roles     RoleMap
}

func (t *Session) touch() {
	t.LastTouch = time.Now()
}

//map of session to user ID
var sessionCache map[string]*Session

func InitSessionCache() error {
	sessionCache = make(map[string]*Session, 0)
	return nil
}
func LookupSession(sessionID string) *Session {
	var ret *Session
	ret = sessionCache[sessionID]
	return ret
}
func RemoveSession(sessionID string) {
	sessionCache[sessionID] = nil
}
func MakeNewSession(userID string) *Session {
	session := Session{UserID: userID, SessionID: uuid.New().String()}
	session.touch()
	user, e := FetchUserUserID(userID)
	if e != nil {
		session.Roles, _ = FetchRoleForUser(user.UserUUID)
	}
	sessionCache[session.SessionID] = &session
	return &session
}

type User struct {
	UserUUID     string `gorm:"type: varchar(36);primary_key"`
	UserID       string `gorm:"type: varchar(255);unique_index;not null"`
	Name         string
	Phone        string
	Email        string
	PasswordHash string `gorm:"type: varchar(255)"`
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
	RoleID RoleType `gorm:"type: varchar(36);primary_key"`
}

type CipherRecord struct {
	IV            string
	RecordID      string
	CipherVersion int
	Data          string
}

func ValidatePassword(userID string, password string) bool {
	user, e := FetchUserUserID(userID)
	if e != nil {
		return false
	}
	hash := HashPassword(password)
	return hash == user.PasswordHash
}

const somedata = "this is a error message"

//recordID is the id of the record, does not matter what it is as long as its a value
//plain is any jsonifiable record
func PlainStructToCipher(recordID string, plain interface{}) *CipherRecord {
	var rec CipherRecord
	rec.RecordID = recordID
	rec.CipherVersion = 1
	hashbits := sha256.Sum256([]byte(somedata + recordID))
	block, err := aes.NewCipher(hashbits[:])
	if err != nil {
		panic(err.Error())
	}

	plainBits, _ := json.Marshal(plain)
	blkSize := block.BlockSize()
	plainBits = PKCS5Padding(plainBits, blkSize)
	ivHash := sha256.Sum256([]byte(time.Now().String()))
	iv := ivHash[0:blkSize]

	cbcEncrypter := cipher.NewCBCEncrypter(block, iv[0:blkSize])
	cipherBits := make([]byte, len(plainBits))
	cbcEncrypter.CryptBlocks(cipherBits, plainBits)
	rec.Data = hex.EncodeToString(cipherBits)
	rec.IV = hex.EncodeToString(iv)

	return &rec

}

func CipherRecordToPlainRecord(cipherIn *CipherRecord, plainOut interface{}) error {
	iv, err := hex.DecodeString(cipherIn.IV)
	if err != nil {
		return err
	}
	hashbits := sha256.Sum256([]byte(somedata + cipherIn.RecordID))
	block, err := aes.NewCipher(hashbits[:])
	if err != nil {
		return err
	}
	cipherBits, err := hex.DecodeString(cipherIn.Data)
	if err != nil {
		return err
	}
	cbcDecrypt := cipher.NewCBCDecrypter(block, iv)
	plainBits := make([]byte, len(cipherBits))
	cbcDecrypt.CryptBlocks(plainBits, cipherBits)
	plainBits = PKCS5Trimming(plainBits)
	err = json.Unmarshal(plainBits, plainOut)
	return err
}

func HashPassword(password string) string {
	hashbits := sha256.Sum256([]byte(password))
	hash := hex.EncodeToString(hashbits[:])
	ret := fmt.Sprintf("%04d%s", 1, hash)
	return ret
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

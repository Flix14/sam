package modelos

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Model
	Username   string    `gorm:"unique;size:250;not null;" sql:"index" json:"username"`
	Name       string    `gorm:"size:250;not null;" json:"name"`
	Password   string    `json:"password,omitempty"`
	FirstLogin time.Time `sql:"index" json:"first_login"`
	LastLogin  time.Time `sql:"index" json:"last_login"`
}

func InitDefaultUser(db *gorm.DB) {
	var password = "1234"
	hashPassword := md5.Sum([]byte(password))
	db.FirstOrCreate(&User{Model: Model{ID: 1}, Username: "temp", Name: "TEMP USER", Password: hex.EncodeToString(hashPassword[:])})
}

func (u *User) Add() (*User, error) {
	hashPassword := md5.Sum([]byte(u.Password))
	u.Password = hex.EncodeToString(hashPassword[:])
	err := DB.Create(&u).Error
	if err != nil {
		return nil, err
	}
	return u, err
}

func (u *User) Find(id int) (err error) {
	err = DB.First(&u, id).Error
	if err != nil {
		return
	}
	return
}

func (u *User) Update() (*User, error) {
	var uu User
	err := DB.First(&uu, u.ID).Error
	if err != nil {
		return nil, err
	}
	if u.Password != "" {
		hashPassword := md5.Sum([]byte(u.Password))
		u.Password = hex.EncodeToString(hashPassword[:])
	}
	err = DB.Save(&u).Error
	if err != nil {
		return nil, err
	}
	return u, err
}

func (u *User) Remove() (err error) {
	err = DB.First(&u, u.ID).Error
	if err != nil {
		return
	}
	err = DB.Delete(&u).Error
	if err != nil {
		return
	}
	return
}

func (u *User) VerifyCredentials(user string, password string) (pu *User, err error) {
	hashPassword := md5.Sum([]byte(password))
	err = DB.Where("username=? AND password=?", user, hex.EncodeToString(hashPassword[:])).Find(&u).Error
	if err != nil {
		return
	}
	pu = u //fix, copy after reassign

	if u.FirstLogin.Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
		err = DB.Model(&u).UpdateColumn("first_login", time.Now()).Error
		if err != nil {
			return
		}
	}
	err = DB.Model(&u).UpdateColumn("last_login", time.Now()).Error
	if err != nil {
		return
	}
	return
}

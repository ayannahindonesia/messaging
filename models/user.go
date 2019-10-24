package models

import (
	"gitlab.com/asira-ayannah/basemodel"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	basemodel.BaseModel
	Username string `json:"username" gorm:"column:username;type:varchar(255);unique;not null"`
	Email    string `json:"email" gorm:"column:email;type:varchar(255);unique;not null"`
	Password string `json:"password" gorm:"column:password;type:text;not null"`
}

func (u *Users) BeforeCreate() (err error) {
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(passwordByte)
	return nil
}

func (i *Users) Create() (err error) {
	err = basemodel.Create(&i)
	return err
}

func (i *Users) Save() (err error) {
	err = basemodel.Save(&i)
	return err
}

func (i *Users) Delete() (err error) {
	err = basemodel.Delete(&i)
	return err
}

func (l *Users) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&l, filter)
	return err
}

func (u *Users) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	users := []Users{}
	order := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.PagedFindFilter(&users, page, rows, order, sorts, filter)

	return result, err
}

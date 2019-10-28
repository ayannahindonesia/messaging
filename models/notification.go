package models

import (
	"time"

	"gitlab.com/asira-ayannah/basemodel"
)

type (
	Notification struct {
		basemodel.BaseModel
		ClientID    int    `json:"client_id" gorm:"column:client_id"`
		Title       string `json:"title" gorm:"column:title"`
		MessageBody string `json:"message_body" gorm:"column:message_body"`
		//TODO: to get from client device
		FirebaseToken string `json:"firebase_token" gorm:"column:firebase_token""`
		//TODO: "direct", "promotion", "failed"; ALLOW NULL
		Topic    string    `json:"topic" gorm:"column:topic"`
		SendTime time.Time `json:"send_time" gorm:"column:send_time" sql:"DEFAULT:current_timestamp"`
	}
)

// gorm callback hook
func (u *Notification) BeforeCreate() (err error) {
	return nil
}

func (u *Notification) Create() error {
	err := basemodel.Create(&u)
	return err
}

// gorm callback hook
func (u *Notification) BeforeSave() (err error) {
	return nil
}

func (u *Notification) Save() error {
	err := basemodel.Save(&u)
	return err
}

func (u *Notification) FindbyID(id int) error {
	err := basemodel.FindbyID(&u, id)
	return err
}

func (u *Notification) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&u, filter)
	return err
}

func (u *Notification) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	notif := []Notification{}
	order := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.PagedFindFilter(&notif, page, rows, order, sorts, filter)

	return result, err
}

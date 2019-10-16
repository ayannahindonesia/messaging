package models

import (
	"time"

	"gitlab.com/asira-ayannah/basemodel"
)

type (
	Messaging struct {
		basemodel.BaseModel
		PhoneNumber string    `json:"phone_number" gorm:"column:phone_number"`
		Message     string    `json:"message" gorm:"column:message"`
		Partner     string    `json:"partner" gorm:"column:partner"`
		Status      bool      `json:"status" gorm:"column:status;type:boolean" sql:"DEFAULT:FALSE"`
		SendTime    time.Time `json:"send_time" gorm:"column:send_time" sql:"DEFAULT:current_timestamp"`
	}
)

// gorm callback hook
func (u *Messaging) BeforeCreate() (err error) {
	return nil
}

func (u *Messaging) Create() error {
	err := basemodel.Create(&u)
	return err
}

// gorm callback hook
func (u *Messaging) BeforeSave() (err error) {
	return nil
}

func (u *Messaging) Save() error {
	err := basemodel.Save(&u)
	return err
}

func (u *Messaging) FindbyID(id int) error {
	err := basemodel.FindbyID(&u, id)
	return err
}

func (u *Messaging) FilterSearchSingle(filter interface{}) error {
	err := basemodel.SingleFindFilter(&u, filter)
	return err
}

func (u *Messaging) PagedFilterSearch(page int, rows int, orderby string, sort string, filter interface{}) (result basemodel.PagedFindResult, err error) {
	mess := []Messaging{}
	order := []string{orderby}
	sorts := []string{sort}
	result, err = basemodel.PagedFindFilter(&mess, page, rows, order, sorts, filter)

	return result, err
}

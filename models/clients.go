package models

import (
	"github.com/ayannahindonesia/basemodel"
	"github.com/google/uuid"
)

type Clients struct {
	basemodel.BaseModel
	Name   string `json:"name" gorm:"column:name"`
	Secret string `json:"secret" gorm:"column:secret"`
	Role   string `json:"role" gorm:"column:role"`
	Key    string `json:"key" gorm:"column:key"`
}

func (i *Clients) BeforeCreate() (err error) {
	if len(i.Secret) < 1 {
		i.Secret = uuid.New().String()
	}
	return nil
}

func (i *Clients) Create() (err error) {
	err = basemodel.Create(&i)
	return err
}

func (i *Clients) Save() (err error) {
	err = basemodel.Save(&i)
	return err
}

func (i *Clients) Delete() (err error) {
	err = basemodel.Delete(&i)
	return err
}

func (l *Clients) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&l, filter)
	return err
}

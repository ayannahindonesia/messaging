package models

import (
	"github.com/google/uuid"
	"gitlab.com/asira-ayannah/basemodel"
)

type Internals struct {
	basemodel.BaseModel
	Name   string `json:"name" gorm:"column:name"`
	Secret string `json:"secret" gorm:"column:secret"`
	Role   string `json:"role" gorm:"column:role"`
	Key    string `json:"key" gorm:"column:key"`
}

func (i *Internals) BeforeCreate() (err error) {
	if len(i.Secret) < 1 {
		i.Secret = uuid.New().String()
	}
	return nil
}

func (i *Internals) Create() (err error) {
	err = basemodel.Create(&i)
	return err
}

func (i *Internals) Save() (err error) {
	err = basemodel.Save(&i)
	return err
}

func (i *Internals) Delete() (err error) {
	err = basemodel.Delete(&i)
	return err
}

func (l *Internals) FilterSearchSingle(filter interface{}) (err error) {
	err = basemodel.SingleFindFilter(&l, filter)
	return err
}

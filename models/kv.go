package models

import (
	"github.com/jinzhu/gorm"
)

type Kv struct {
	gorm.Model
	Key string
	Val []byte `gorm:"type:text(0)"`
}

func SetKv(k string, v []byte) error {
	kv := Kv{
		Key: k,
		Val: v,
	}
	return DB.Create(&kv).Error
}

func GetKv(k string) ([]byte, error) {
	var vv Kv
	err := DB.Where(&Kv{Key: k}).First(&vv).Error
	if err != nil {
		return nil, err
	}
	return vv.Val, nil
}

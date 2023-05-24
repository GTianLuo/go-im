package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	NickName string `gorm:"type:varchar(25) not null"`
	Account  string `gorm:"type:varchar(20) not null unique"`
	Password string `gorm:"type:varchar(255)"`
}

func (*User) TableName() string {
	return "user"
}

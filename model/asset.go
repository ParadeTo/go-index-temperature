package model

type Asset struct {
	Code string `gorm:"type:varchar(10);primary_key"`
	Date string `gorm:"size:10;primary_key"`
	Price float32
	Cap float64
	Pe float32
	Pb float32
}
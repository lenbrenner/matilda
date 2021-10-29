package main

import "gorm.io/gorm"

type Person struct {
	gorm.Model
	Name  string
	Speed float32
}

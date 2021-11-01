package main

import "gorm.io/gorm"

type Person2 struct {
	gorm.Model
	Name  string
	Speed float32
}

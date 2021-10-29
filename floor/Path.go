package floor

import (
	"gorm.io/gorm"
)

type Path struct {
	gorm.Model
	start Location
	end   Location
	nodes []Location
}

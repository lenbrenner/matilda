package floor

import (
	"encoding/json"
	"fmt"
	"takeoff.com/matilda/util"
)

type Plan struct {
	Locations []Location
}

func (plan *Plan) LoadJson(filename string) error {
	dat, err := util.Resource(fmt.Sprintf("maps/%s.json", filename))
	if err == nil {
		err := json.Unmarshal(dat, &plan)
		return err
	} else {
		return err
	}
}

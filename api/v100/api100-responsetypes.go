package api100

import "github.com/chaosvermittlung/funkloch-server/db/v100"

type storeItemCountResponse struct {
	Name  string
	Count int
}

type equipmentCountResponse struct {
	Equipment db100.Equipment
	Store     db100.Store
	Count     int
}

package api100

import (
	"time"

	"github.com/Chaosvermittlung/funkloch-server/db/v100"
)

type authResponse struct {
	Token string `json:"token"`
}

type storeItemCountResponse struct {
	Name  string
	Count int
}

type itemResponse struct {
	Item      db100.Item
	Store     db100.Store
	Box       db100.Box
	Equipment db100.Equipment
}

type equipmentCountResponse struct {
	Equipment db100.Equipment
	Store     db100.Store
	Count     int
}

type eventParticipiantsResponse struct {
	User      db100.User
	Arrival   time.Time
	Departure time.Time
}

type wishlistItemsResponse struct {
	Equipment db100.Equipment
	Count     int
}

type packinglistItemsResponse struct {
	StoreItemID int
	Equipment   db100.Equipment
	Store       db100.Store
}

type faultResponse struct {
	Fault db100.Fault
	Code  int
	Name  string
}

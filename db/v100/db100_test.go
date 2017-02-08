package db100

import (
	"log"
	"os"
	"testing"

	"time"

	"github.com/chaosvermittlung/funkloch-server/global"
)

func TestMain(m *testing.M) {
	var con global.DBConnection
	con.Driver = "sqlite3"
	con.Connection = "./test.db"
	os.Remove(con.Connection)

	Initialisation(&con)
	exit := m.Run()

	err := os.Remove(con.Connection)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exit)
}

func TestUserInsert(t *testing.T) {

	u := User{UserID: -1, Username: "test", Password: "test", Email: "test@test", Right: USERRIGHT_ADMIN}
	s, err := global.GenerateSalt()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	u.Salt = s

	pw, err := global.GeneratePasswordHash(u.Password, u.Salt)
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	u.Password = pw

	err = u.Insert()

	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}

	if u.UserID < 0 {
		t.Errorf("Expected Userid > 0 but got %v", u.UserID)
	}
}

func TestDoesUserExist(t *testing.T) {
	cont, err := DoesUserExist("foobar")

	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}

	if cont {
		t.Errorf("Expected false got %v", cont)
	}

	cont, err = DoesUserExist("test")

	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}

	if !cont {
		t.Errorf("Expected true got %v", cont)
	}
}

func TestGetUsers(t *testing.T) {
	uu, err := GetUsers()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if len(uu) != 2 {
		t.Error("Expected length 1 got ", len(uu))
	}
}

func TestGetDetailstoUsername(t *testing.T) {
	u := User{Username: "test"}
	err := u.GetDetailstoUsername()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if u.UserID != 2 {
		t.Errorf("Expected User_ID 1 but got %v", u.UserID)
	}
}

func TestGetUserDetails(t *testing.T) {
	u := User{UserID: 2}
	err := u.GetDetails()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if u.Username != "test" {
		t.Errorf("Expected Username netfabb but got %v", u.Username)
	}
}

func TestPatchUser(t *testing.T) {
	u := User{UserID: -1, Username: "test", Password: "test", Email: "test@test", Right: USERRIGHT_ADMIN}
	un := User{UserID: -1, Username: "test", Password: "test", Email: "test@otherhost", Right: USERRIGHT_MEMBER}

	err := u.Patch(un)

	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}

	if u.Email != "test@otherhost" {
		t.Errorf("Expected Email admin@otherhost but got %v", u.Email)
	}

	if u.Right != USERRIGHT_MEMBER {
		t.Errorf("Expected Right Operate but got %v", u.Right)
	}
}

func TestUpdateUser(t *testing.T) {
	u := User{UserID: 2, Username: "test1", Password: "test", Email: "test@otherhost", Right: USERRIGHT_ADMIN}
	un := User{UserID: 2}
	err := u.Update()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	err = un.GetDetails()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if u.Username != un.Username {
		t.Errorf("Found differences between old and updated username: %v %v", u.Username, un.Username)
	}
	if u.Email != un.Email {
		t.Errorf("Found differences between old and updated username: %v %v", u.Email, un.Email)
	}
}

func TestDeleteUser(t *testing.T) {
	err := DeleteUser(2)
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	uu, err := GetUsers()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if len(uu) != 1 {
		t.Error("Expected length 0 got ", len(uu))
	}
}

func TestStoreInsert(t *testing.T) {
	s := Store{StoreID: -1, Name: "foobar", Adress: "test", Manager: 1}
	err := s.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}

	if s.StoreID < 0 {
		t.Errorf("Expected Storeid > 0 but got %v", s.StoreID)
	}
}

func TestGetStores(t *testing.T) {
	ss, err := GetStores()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(ss) != 1 {
		t.Errorf("Expected len = 1 got %v", len(ss))
	}
}

func TestStoreGetDetails(t *testing.T) {
	s := Store{StoreID: 1}
	err := s.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if s.Name != "foobar" {
		t.Errorf("Expected Name = foobar but got %v", s.Name)
	}
	if s.Manager != 1 {
		t.Errorf("Expected Manager = 1 but got %v", s.Manager)
	}
	if s.Adress != "test" {
		t.Errorf("Expected Adress = test but got %v", s.Adress)
	}
}

func TestStoreUpdate(t *testing.T) {
	s := Store{StoreID: 1, Name: "foobar2", Adress: "test2", Manager: 1}
	sn := Store{StoreID: 1}
	err := s.Update()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = sn.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if sn.Name != s.Name {
		t.Error("Name missmatch:", sn.Name, s.Name)
	}
	if sn.Adress != s.Adress {
		t.Error("Adress missmatch:", sn.Adress, s.Adress)
	}
	if sn.Manager != s.Manager {
		t.Error("Manager missmatch:", sn.Manager, s.Manager)
	}
}

func TestStoreGetManager(t *testing.T) {
	s := Store{StoreID: 1}
	err := s.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	m, err := s.GetManager()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if m.UserID != 1 {
		t.Errorf("Epected UserId = 1 but got %v", m.UserID)
	}
}

func TestStoreDelete(t *testing.T) {
	s := Store{StoreID: 1}
	err := s.Delete()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestEquipmentInsert(t *testing.T) {
	e := Equipment{EquipmentID: -1, Name: "FF54"}
	err := e.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if e.EquipmentID != 1 {
		t.Errorf("Expected EquipmentID = 1 but got %v", e.EquipmentID)
	}
}

func TestEquipmentGetDetails(t *testing.T) {
	e := Equipment{EquipmentID: 1}
	err := e.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if e.Name != "FF54" {
		t.Error("Expected Name = FF54 but got", e.Name)
	}
}

func TestEquipmentUpdate(t *testing.T) {
	e := Equipment{EquipmentID: 1, Name: "FF33"}
	en := Equipment{EquipmentID: 1}
	err := e.Update()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = en.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if e.Name != en.Name {
		t.Error("Name missmatch:", e.Name, en.Name)
	}
}

func TestEquipmentDelete(t *testing.T) {
	e := Equipment{EquipmentID: 1}
	err := e.Delete()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestStoreItemInsert(t *testing.T) {
	e := Equipment{EquipmentID: -1, Name: "FF54"}
	s := Store{StoreID: -1, Name: "foobar", Adress: "test", Manager: 1}
	err := e.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = s.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	si := StoreItem{StoreItemID: -1, StoreID: s.StoreID, EquipmentID: e.EquipmentID}
	err = si.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if si.StoreItemID != 1 {
		t.Errorf("Expected StoreItemID = 1 but got %v", si.StoreItemID)
	}
}

func TestStoreItemGetDetails(t *testing.T) {
	si := StoreItem{StoreItemID: 1}
	err := si.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if si.StoreID != 2 {
		t.Errorf("Expected StoreID = 2 but got %v", si.StoreID)
	}
	if si.EquipmentID != 2 {
		t.Errorf("Expected EquipmentID = 2 but got %v", si.EquipmentID)
	}
}

func TestStoreItemUpdate(t *testing.T) {
	s := Store{StoreID: -1, Name: "foobar", Adress: "test", Manager: 1}
	err := s.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	si := StoreItem{StoreItemID: 1, StoreID: s.StoreID}
	sin := StoreItem{StoreItemID: 1}
	err = si.Update()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = sin.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if si.StoreID != sin.StoreID {
		t.Errorf("StoreID missmatch: %v %v", si.StoreID, sin.StoreID)
	}
}

func TestStoreItemPostFault(t *testing.T) {
	si := StoreItem{StoreItemID: 1}
	f := Fault{FaultID: -1, StoreItemID: si.StoreItemID, Status: FaultStatusNew, Comment: "Alles kaputt"}
	f, err := si.PostFault(f)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if f.FaultID != 1 {
		t.Errorf("Expected FaultID = 1 but got %v", f.FaultID)
	}
}

func TestStoreItemGetFaults(t *testing.T) {
	si := StoreItem{StoreItemID: 1}
	ff, err := si.GetFaults()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(ff) != 1 {
		t.Errorf("Expected len = 1 but got %v", len(ff))
	}
}

func TestStoreGetStoreItems(t *testing.T) {
	s := Store{StoreID: 3}
	sis, err := s.GetStoreitems()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(sis) != 1 {
		t.Errorf("Expected len = 1 but got %v", len(sis))
	}
}

func TestStoreGetItemCount(t *testing.T) {
	s := Store{StoreID: 3}
	c, err := s.GetItemCount(2)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if c != 1 {
		t.Errorf("Expected len = 1 but got %v", c)
	}
}

func TestGetFaults(t *testing.T) {
	ff, err := GetFaults()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(ff) != 1 {
		t.Errorf("Expected len = 1 but got %v", len(ff))
	}
}

func TestGetFaultDetails(t *testing.T) {
	f := Fault{FaultID: 1}
	err := f.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if f.Status != FaultStatusNew {
		t.Errorf("Expected Status New but got %v", f.Status)
	}
	if f.Comment != "Alles kaputt" {
		t.Errorf("Expected Comment Alles kaputt but got %v", f.Comment)
	}
}

func TestUpdateFault(t *testing.T) {
	f := Fault{FaultID: 1, Status: FaultStatusFixed, Comment: "Nix mehr kaputt"}
	fn := Fault{FaultID: 1}
	err := f.Update()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = fn.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if f.Status != fn.Status {
		t.Errorf("Status missmatch: %v %v", f.Status, fn.Status)
	}
	if f.Comment != fn.Comment {
		t.Errorf("Comment missmatch: %v %v", f.Comment, fn.Comment)
	}
}

func TestDeleteFault(t *testing.T) {
	f := Fault{FaultID: 1}
	err := f.Delete()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestStoreItemDelete(t *testing.T) {
	si := StoreItem{StoreID: 1}
	err := si.Delete()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestEventInsert(t *testing.T) {
	e := Event{EventID: -1, Name: "CCS", Adress: "Chiemsee", Start: time.Now().Add(time.Hour * 24), End: time.Now().Add(time.Hour * 24 * 3)}
	err := e.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if e.EventID != 1 {
		t.Errorf("Expected EventID = 1 but got %v", e.EventID)
	}
}

func TestEventGetDetails(t *testing.T) {
	e := Event{EventID: 1}
	err := e.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if e.Name != "CCS" {
		t.Error("Expected name CCS but got", e.Name)
	}
	if e.Adress != "Chiemsee" {
		t.Error("Expected adress Chiemsee but got", e.Adress)
	}
}

func TestEventUpdate(t *testing.T) {
	e := Event{EventID: 1, Name: "CSS", Adress: "Würmsee", Start: time.Now().Add(time.Hour * 24), End: time.Now().Add(time.Hour * 24 * 3)}
	en := Event{EventID: 1}
	err := e.Update()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = en.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if e.Name != en.Name {
		t.Error("Name missmatch:", e.Name, en.Name)
	}
	if e.Adress != en.Adress {
		t.Error("Adress missmatch:", e.Adress, en.Adress)
	}
	if !e.Start.Equal(en.Start) {
		t.Errorf("Start missmatch: %v %v", e.Start, en.Start)
	}
	if !e.End.Equal(en.End) {
		t.Errorf("End missmatch: %v %v", e.End, en.End)
	}
}

func TestGetEvents(t *testing.T) {
	ee, err := GetEvents()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(ee) != 1 {
		t.Errorf("Expected len = 1 but bot %v", len(ee))
	}
}

func TestGetNextEvent(t *testing.T) {
	e, err := GetNextEvent()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if e.EventID != 1 {
		t.Errorf("Expected EventID = 1 but got %v", e.EventID)
	}
}

func TestParticipantInsert(t *testing.T) {
	p := Participant{UserID: 1, EventID: 1, Arrival: time.Now().Add(25 * time.Hour), Departure: time.Now().Add(23 * time.Hour)}
	err := p.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestEventGetParticipants(t *testing.T) {
	e := Event{EventID: 1}
	pp, err := e.GetParticipants()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(pp) != 1 {
		t.Errorf("Expected len = 1 but got %v", len(pp))
	}
	if pp[0].UserID != 1 {
		t.Errorf("Excpected UserID = 1 but got %v", pp[0].UserID)
	}
}

func TestParticipantGetDetails(t *testing.T) {
	p := Participant{UserID: 1, EventID: 1}
	err := p.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestParticipantUpdate(t *testing.T) {
	p := Participant{UserID: 1, EventID: 1, Arrival: time.Now().Add(25 * time.Hour), Departure: time.Now().Add(23 * time.Hour)}
	pn := Participant{UserID: 1, EventID: 1}
	err := p.Update()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = pn.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if !p.Arrival.Equal(pn.Arrival) {
		t.Errorf("Arrival missmatch: %v %v", p.Arrival, pn.Arrival)
	}
	if !p.Departure.Equal(pn.Departure) {
		t.Errorf("Departure missmatch: %v %v", p.Departure, pn.Departure)
	}
}

func TestParticipantDelete(t *testing.T) {
	p := Participant{UserID: 1, EventID: 1}
	err := p.Delete()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestEventDelete(t *testing.T) {
	e := Event{EventID: 1}
	err := e.Delete()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestPackinglistInsert(t *testing.T) {
	e := Event{EventID: 1, Name: "CSS", Adress: "Würmsee", Start: time.Now().Add(time.Hour * 24), End: time.Now().Add(time.Hour * 24 * 3)}
	err := e.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	pl := Packinglist{PackinglistID: -1, Name: "Galaxy", EventID: 2}
	err = pl.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if pl.PackinglistID != 1 {
		t.Errorf("Expected PackinglistID = 1 but got %v", pl.PackinglistID)
	}
}

func TestGetPackingLists(t *testing.T) {
	pp, err := GetPackinglists()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(pp) != 1 {
		t.Errorf("Expected len = 1 but got %v", len(pp))
	}
}

func TestPackingListGetDetails(t *testing.T) {
	p := Packinglist{PackinglistID: 1}
	err := p.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if p.Name != "Galaxy" {
		t.Error("Expected Name = Galaxy but got", p.Name)
	}
}

func TestPackingListUpdate(t *testing.T) {
	p := Packinglist{PackinglistID: 1, Name: "Focus", EventID: 2}
	pn := Packinglist{PackinglistID: 1}
	err := p.Update()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = pn.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if p.Name != pn.Name {
		t.Error("Name missmatch:", p.Name, pn.Name)
	}
}

func TestPackingListItemInsert(t *testing.T) {
	e := Equipment{Name: "FF54"}
	err := e.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	s := Store{Name: "Kitchen", Manager: 1}
	err = s.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	si := StoreItem{StoreID: s.StoreID, EquipmentID: e.EquipmentID}
	err = si.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	pli := PackinglistItem{PackinglistID: 1, StoreitemID: si.StoreItemID}
	err = pli.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestPackinglistGetItems(t *testing.T) {
	p := Packinglist{PackinglistID: 1}
	sis, err := p.GetItems()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(sis) != 1 {
		t.Errorf("Expected len = 1 but got %v", err)
	}
}

func TestPackinglistItemDelete(t *testing.T) {
	p := PackinglistItem{StoreitemID: 2, PackinglistID: 1}
	err := p.Delete()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestPackingListDelete(t *testing.T) {
	p := Packinglist{PackinglistID: 1}
	err := p.Delete()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestWishlistInsert(t *testing.T) {
	w := Wishlist{WishlistID: -1, Name: "test"}
	err := w.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if w.WishlistID != 1 {
		t.Errorf("Expected WishlistID = 1 but got %v", w.WishlistID)
	}
}

func TestGetWishlists(t *testing.T) {
	ww, err := GetWishlists()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(ww) != 1 {
		t.Errorf("Epected len = 1 but got %v", len(ww))
	}
}

func TestWishlistGetDetails(t *testing.T) {
	w := Wishlist{WishlistID: 1}
	err := w.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if w.Name != "test" {
		t.Error("Expected Name = test but got", w.Name)
	}
}

func TestWishlistUpdate(t *testing.T) {
	w := Wishlist{WishlistID: 1, Name: "test2"}
	wn := Wishlist{WishlistID: 1}
	err := w.Update()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = wn.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if w.Name != wn.Name {
		t.Error("Name missmatch:", w.Name, wn.Name)
	}
}

func TestWishlistItemInsert(t *testing.T) {
	wli := Wishlistitem{WishlistID: 1, EquipmentID: 3, Count: 5}
	err := wli.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestWishlistItemGetDetails(t *testing.T) {
	wli := Wishlistitem{WishlistID: 1, EquipmentID: 3}
	err := wli.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if wli.Count != 5 {
		t.Errorf("Expected Count = 5 but got %v", wli.Count)
	}
}

func TestWishlistGetItems(t *testing.T) {
	w := Wishlist{WishlistID: 1}
	ee, err := w.GetItems()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(ee) != 1 {
		t.Errorf("Expected len = 1 but got %v", len(ee))
	}
}

func TestWishlistItemUpdate(t *testing.T) {
	wli := Wishlistitem{WishlistID: 1, EquipmentID: 3, Count: 10}
	wlin := Wishlistitem{WishlistID: 1, EquipmentID: 3}
	err := wli.Update()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = wlin.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if wli.Count != wlin.Count {
		t.Errorf("Count missmatch: %v %v", wli.Count, wlin.Count)
	}
}

func TestWishlistItemDelete(t *testing.T) {
	wli := Wishlistitem{WishlistID: 1, EquipmentID: 3}
	err := wli.Delete()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestWishlistDelete(t *testing.T) {
	wl := Wishlist{WishlistID: 1}
	err := wl.Delete()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

package db100

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"time"

	"github.com/Chaosvermittlung/funkloch-server/global"
)

func TestMain(m *testing.M) {
	var con global.DBConnection
	con.Driver = "sqlite3"
	con.Connection = "./test.db"
	os.Remove(con.Connection)
	abs, _ := filepath.Abs(con.Connection)
	log.Println("Test Database Path:", abs)
	Initialisation(&con)
	exit := m.Run()

	/*err := os.Remove(con.Connection)
	if err != nil {
		log.Fatal(err)
	}*/
	os.Exit(exit)
}

func TestUserInsert(t *testing.T) {

	u := User{Username: "test", Password: "test", Email: "test@test", Right: USERRIGHT_ADMIN}
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
	u := User{Username: "test", Password: "test", Email: "test@test", Right: USERRIGHT_ADMIN}
	un := User{Username: "test", Password: "test", Email: "test@otherhost", Right: USERRIGHT_MEMBER}

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
	var m User
	m.UserID = 1
	err := m.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	s := Store{Name: "foobar", Adress: "test", Manager: m, ManagerID: 1}
	err = s.Insert()
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
	if s.ManagerID != 1 {
		t.Errorf("Expected Manager = 1 but got %v", s.Manager)
	}
	if s.Adress != "test" {
		t.Errorf("Expected Adress = test but got %v", s.Adress)
	}
}

func TestStoreUpdate(t *testing.T) {
	s := Store{StoreID: 1, Name: "foobar2", Adress: "test2", ManagerID: 1}
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
	e := Equipment{Name: "FF54"}
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

func TestGetEquipment(t *testing.T) {
	res, err := GetEquipment()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(res) < 1 {
		t.Errorf("Expected len > 1 got %v", len(res))
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

func TestBoxInsert(t *testing.T) {
	s := Store{Name: "foobar", Adress: "test", ManagerID: 1}
	err := s.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	b := Box{StoreID: s.StoreID, Description: "TestBox"}
	err = b.Insert()
	if b.BoxID != 1 {
		t.Errorf("Expected BoxID = 1 but got %v", b.BoxID)
	}
	if b.Code != 2020000000013 {
		t.Errorf("Expected Code = 2020000000013 but got %v", b.Code)
	}
}

func TestGetBoxes(t *testing.T) {
	bb, err := GetBoxes()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(bb) != 1 {
		t.Errorf("Expected len = 1 got %v", len(bb))
	}
}

func TestBoxGetDetails(t *testing.T) {
	b := Box{BoxID: 1}
	err := b.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if b.StoreID != 2 {
		t.Error("Expected StoreID = 2 but got", b.StoreID)
	}
	if b.Code != 2020000000013 {
		t.Error("Expected Code = 2020000000013 but got", b.StoreID)
	}
	if b.Description != "TestBox" {
		t.Error("Expected Name = TestBox but got", b.Description)
	}
}

func TestBoxUpdate(t *testing.T) {
	b := Box{BoxID: 1, StoreID: 2, Description: "TestBox2"}
	bn := Box{BoxID: 1}
	err := b.Update()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = bn.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if bn.Description != "TestBox2" {
		t.Error("Expected Name = TestBox2 but got", bn.Description)
	}
}

func TestStoreAddStoreBox(t *testing.T) {
	b := Box{BoxID: 1}
	err := b.GetDetails()
	if err != nil {
		t.Errorf("Expected no #1 error but got %v", err)
	}
	s := Store{StoreID: 2}
	err = s.GetDetails()
	if err != nil {
		t.Errorf("Expected no #2 error but got %v", err)
	}
	err = s.AddStoreBox(b)
	if err != nil {
		t.Errorf("Expected no #3 error but got %v", err)
	}
}

func TestStoreGetStoreBoxes(t *testing.T) {
	s := Store{StoreID: 2}
	err := s.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	bb, err := s.GetStoreBoxes()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(bb) != 1 {
		t.Errorf("Expected len = 1 got %v", len(bb))
	}
}

func TestBoxDelete(t *testing.T) {
	b := Box{BoxID: 1}
	err := b.Delete()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestItemInsert(t *testing.T) {
	e := Equipment{Name: "FF54"}
	s := Store{Name: "foobar", Adress: "test", ManagerID: 1}
	err := e.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = s.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	b := Box{StoreID: s.StoreID, Description: "TestBox"}
	err = b.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	si := Item{BoxID: b.BoxID, EquipmentID: e.EquipmentID, Description: "Foobar"}
	err = si.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if si.ItemID != 1 {
		t.Errorf("Expected ItemID = 1 but got %v", si.ItemID)
	}
	if si.Code != 2000000000015 {
		t.Errorf("Expected Code = 2000000000015 but got %v", si.Code)
	}
}

func TestItemGetDetails(t *testing.T) {
	si := Item{ItemID: 1}
	err := si.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if si.BoxID != 2 {
		t.Errorf("Expected BoxID = 2 but got %v", si.BoxID)
	}
	if si.EquipmentID != 2 {
		t.Errorf("Expected EquipmentID = 2 but got %v", si.EquipmentID)
	}
	if si.Code != 2000000000015 {
		t.Errorf("Expected Code = 2000000000015 but got %v", si.Code)
	}
	if si.Description != "Foobar" {
		t.Errorf("Expected Description = Foobar but got %v", si.Description)
	}
}

func TestBoxAddBoxItem(t *testing.T) {
	b := Box{BoxID: 2}
	si := Item{ItemID: 1}
	err := b.AddBoxItem(si)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestBoxGetBoxItems(t *testing.T) {
	b := Box{BoxID: 2}
	res, err := b.GetBoxItems()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(res) != 1 {
		t.Errorf("Expected len = 1 got %v", len(res))
	}
}

func TestItemUpdate(t *testing.T) {
	e := Equipment{Name: "FF54"}
	b := Box{StoreID: 2, Description: "TestBox2"}
	err := b.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = e.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	si := Item{ItemID: 1, BoxID: b.BoxID, EquipmentID: e.EquipmentID}
	sin := Item{ItemID: 1}
	err = si.Update()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = sin.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if si.BoxID != sin.BoxID {
		t.Errorf("StoreID missmatch: %v %v", si.BoxID, sin.BoxID)
	}
	if si.EquipmentID != sin.EquipmentID {
		t.Errorf("EquipmentID missmatch: %v %v", si.EquipmentID, sin.EquipmentID)
	}
}

func TestGetItems(t *testing.T) {
	ii, err := GetItems(false)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(ii) < 1 {
		t.Errorf("Expected len > 1 got %v", len(ii))
	}

	ii, err = GetItems(true)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(ii) != 0 {
		t.Errorf("Expected len > 1 got %v", len(ii))
	}
}

func TestItemAddFault(t *testing.T) {
	si := Item{ItemID: 1}
	f := Fault{ItemID: si.ItemID, Status: FaultStatusNew, Comment: "Alles kaputt"}
	f, err := si.AddFault(f)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if f.FaultID != 1 {
		t.Errorf("Expected FaultID = 1 but got %v", f.FaultID)
	}
}

func TestItemGetFaults(t *testing.T) {
	si := Item{ItemID: 1}
	ff, err := si.GetFaults()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(ff) != 1 {
		t.Errorf("Expected len = 1 but got %v", len(ff))
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

func TestItemDelete(t *testing.T) {
	si := Item{ItemID: 1}
	err := si.Delete()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestEventInsert(t *testing.T) {
	e := Event{Name: "CCS", Adress: "Chiemsee", Start: time.Now().Add(time.Hour * 24), End: time.Now().Add(time.Hour * 24 * 3)}
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
	pl := Packinglist{Name: "Galaxy", EventID: 1}
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

func TestEventGetPackinglists(t *testing.T) {
	e := Event{EventID: 1}
	pp, err := e.GetPackinglists()
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
	p := Packinglist{PackinglistID: 1, Name: "Focus", EventID: 1}
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

func TestPackinglistAddBox(t *testing.T) {
	e := Equipment{Name: "FF54"}
	err := e.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	s := Store{Name: "Kitchen", Adress: "Straße", ManagerID: 1}
	err = s.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	b := Box{StoreID: s.StoreID, Description: "TestBox2"}
	err = b.Insert()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	p := Packinglist{PackinglistID: 1}
	err = p.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	err = p.AddPackinglistBox(b)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestPackinglistGetPackinglistBoxes(t *testing.T) {
	p := Packinglist{PackinglistID: 1}
	sis, err := p.GetPackinglistBoxes()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(sis) != 1 {
		t.Errorf("Expected len = 1 but got %v", err)
	}
}

func TestPackinglistRemovePackinglistBox(t *testing.T) {
	b := Box{BoxID: 3}
	err := b.GetDetails()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	p := Packinglist{PackinglistID: 1}
	err = p.RemovePackinglistBox(b)
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
	w := Wishlist{Name: "test"}
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

func TestWishlistAddWishlistItems(t *testing.T) {
	w := Wishlist{WishlistID: 1}
	e := Equipment{EquipmentID: 4}
	err := w.AddWishlistItem(e)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

func TestWishlistGetWishlistItems(t *testing.T) {
	w := Wishlist{WishlistID: 1}
	ee, err := w.GetWishlistItems()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
	if len(ee) != 1 {
		t.Errorf("Expected len = 1 but got %v", len(ee))
	}
}

func TestWishlistDelete(t *testing.T) {
	wl := Wishlist{WishlistID: 1}
	err := wl.Delete()
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}
}

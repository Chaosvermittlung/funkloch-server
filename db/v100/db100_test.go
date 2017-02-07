package db100

import (
	"log"
	"os"
	"testing"

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

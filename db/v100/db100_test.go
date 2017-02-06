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

	u := User{UserID: -1, Username: "admin", Password: "admin", Email: "admin@localhost", Right: USERRIGHT_ADMIN}
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

	cont, err = DoesUserExist("admin")

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
	if len(uu) != 1 {
		t.Error("Expected length 1 got %v", len(uu))
	}
}

func TestGetDetailstoUsername(t *testing.T) {
	u := User{Username: "admin"}
	err := u.GetDetailstoUsername()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if u.UserID != 1 {
		t.Errorf("Expected User_ID 1 but got %v", u.UserID)
	}
}

func TestGetUserDetails(t *testing.T) {
	u := User{UserID: 1}
	err := u.GetDetails()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if u.Username != "admin" {
		t.Errorf("Expected Username netfabb but got %v", u.Username)
	}
}

func TestPatchUser(t *testing.T) {
	u := User{UserID: -1, Username: "admin", Password: "admin", Email: "admin@localhost", Right: USERRIGHT_ADMIN}
	un := User{UserID: -1, Username: "admin", Password: "admin", Email: "admin@otherhost", Right: USERRIGHT_MEMBER}

	err := u.Patch(un)

	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}

	if u.Email != "admin@otherhost" {
		t.Errorf("Expected Email admin@otherhost but got %v", u.Email)
	}

	if u.Right != USERRIGHT_MEMBER {
		t.Errorf("Expected Right Operate but got %v", u.Right)
	}
}

func TestUpdateUser(t *testing.T) {
	u := User{UserID: 1, Username: "admin1", Password: "admin", Email: "admin@otherhost", Right: USERRIGHT_ADMIN}
	un := User{UserID: 1}
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
	err := DeleteUser(1)
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	uu, err := GetUsers()
	if err != nil {
		t.Fatalf("No error expected but got %v", err)
	}
	if len(uu) != 0 {
		t.Error("Expected length 0 got %v", len(uu))
	}
}

package db100

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chaosvermittlung/funkloch-server/global"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func Initialisation(dbc *global.DBConnection) {
	var err error
	db, err = sqlx.Open(dbc.Driver, dbc.Connection)
	if err != nil {
		log.Fatal(err)
	}
	initDB(dbc)
}

func initDB(dbc *global.DBConnection) {
	switch dbc.Driver {
	case "sqlite3":
		cont, err := global.Exists(dbc.Connection)
		if err != nil {
			log.Fatal(err)
		}
		if cont {
			fmt.Println("cont")
			return
		}
		_, err = os.Create(dbc.Connection)
		if err != nil {
			log.Fatal("Could not create file "+dbc.Connection, err)
		}
		_, err = db.Exec(createSQLlitestmt)
		if err != nil {
			log.Printf("%q: %s\n", err, createSQLlitestmt)
			return
		}
		var u User
		u.Username = "admin"
		u.Password = "admin"
		u.Email = "admin@localhost"
		u.Right = USERRIGHT_ADMIN
		s, err := global.GenerateSalt()
		if err != nil {
			log.Println(err)
		}
		u.Salt = s

		pw, err := global.GeneratePasswordHash(u.Password, u.Salt)
		if err != nil {
			log.Println(err)
		}
		u.Password = pw
		err = u.Insert()
		if err != nil {
			log.Println(err)
		}
	default:
		log.Fatal("DB Driver unkown. Stopping Server")
	}
}

type UserRight int

const (
	USERRIGHT_MEMBER UserRight = 1 + iota
	USERRIGHT_ADMIN
)

type User struct {
	UserID   int       `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Salt     string    `json:"-"`
	Email    string    `json:"email"`
	Right    UserRight `json:"userright"`
}

func copyifnotempty(str1, str2 string) string {
	if str2 != "" {
		return str2
	} else {
		return str1
	}
}

func DoesUserExist(username string) (bool, error) {
	var id int
	err := db.Get(&id, "Select Count(*) from User Where Username = ?", username)
	b := (id > 0)
	return b, err
}

func GetUsers() ([]User, error) {
	var u []User
	err := db.Select(&u, "Select * from User")
	return u, err
}

func (u *User) GetDetailstoUsername() error {
	err := db.Get(u, "SELECT * from User Where Username = ? Limit 1", u.Username)
	return err
}

func (u *User) GetDetails() error {
	err := db.Get(u, "SELECT * from User Where UserID = ? Limit 1", u.UserID)
	return err
}

func (u *User) Patch(ou User) error {
	u.Username = copyifnotempty(u.Username, ou.Username)
	if ou.Password != "" {
		p, err := global.GeneratePasswordHash(ou.Password, u.Salt)
		if err != nil {
			return err
		}
		u.Password = p
	}
	u.Email = copyifnotempty(u.Email, ou.Email)
	if ou.Right != 0 {
		u.Right = ou.Right
	}
	return nil
}

func (u *User) Update() error {
	_, err := db.Exec("UPDATE User SET username = ?, password = ?, email = ?, right = ? WHERE Userid = ?", u.Username, u.Password, u.Email, u.Right, u.UserID)
	return err
}

func (u *User) Insert() error {
	res, err := db.Exec("INSERT INTO User (username, password, salt, email, right) VALUES(?,?,?,?,?)", u.Username, u.Password, u.Salt, u.Email, u.Right)
	if err != nil {
		log.Println(err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	u.UserID = int(id)

	return nil
}

func DeleteUser(id int) error {
	_, err := db.Exec("DELETE FROM User Where UserID = ?", id)
	return err
}

type Store struct {
	StoreID int
	Name    string
	Adress  string
	Manager int
}

func (s *Store) Insert() error {
	res, err := db.Exec("Insert Into Store (name, adress, manager) Values (?,?,?)", s.Name, s.Adress, s.Manager)
	if err != nil {
		log.Println(err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	s.StoreID = int(id)
	return nil
}

func GetStores() ([]Store, error) {
	var s []Store
	err := db.Select(&s, "Select * from Store")
	return s, err
}

func (s *Store) GetDetails() error {
	err := db.Get(s, "SELECT * from Store Where StoreID = ? Limit 1", s.StoreID)
	return err
}

func (s *Store) GetManager() (User, error) {
	var u User
	u.UserID = s.Manager
	err := u.GetDetails()
	return u, err
}

func (s *Store) Update() error {
	_, err := db.Exec("Update Store SET name = ?, adress = ?, manager = ? where StoreID = ?", s.Name, s.Adress, s.Manager, s.StoreID)
	return err
}

func (s *Store) Delete() error {
	_, err := db.Exec("Delete from Store Where StoreID = ?", s.StoreID)
	return err
}

func (s *Store) GetStoreitems() ([]StoreItem, error) {
	var si []StoreItem
	err := db.Select(&si, "Select * from StoreItem Where StoreID = ?", s.StoreID)
	return si, err
}

func (s *Store) GetItemCount(id int) (int, error) {
	var i int
	err := db.Get(&i, "Select Count(*) from StoreItem Where StoreID = ? and EquipmentID = ?", s.StoreID, id)
	return i, err
}

type Equipment struct {
	EquipmentID int
	Name        string
}

func (e *Equipment) Insert() error {
	res, err := db.Exec("Insert Into Equipment (Name) Values (?)", e.Name)
	if err != nil {
		log.Println(err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	e.EquipmentID = int(id)
	return nil
}

func GetEquipment() ([]Equipment, error) {
	var e []Equipment
	err := db.Select(&e, "Select * from Equipment")
	return e, err
}

func (e *Equipment) GetDetails() error {
	err := db.Get(e, "SELECT * from Equipment Where EquipmentID = ? Limit 1", e.EquipmentID)
	return err
}

func (e *Equipment) Update() error {
	_, err := db.Exec("Update Equipment SET name = ? where EquipmentID = ?", e.Name, e.EquipmentID)
	return err
}

func (e *Equipment) Delete() error {
	_, err := db.Exec("Delete from Equipment Where EquipmentID = ?", e.EquipmentID)
	return err
}

type StoreItem struct {
	StoreItemID int
	StoreID     int
	EquipmentID int
}

func (s *StoreItem) Insert() error {
	res, err := db.Exec("Insert Into Sotreitem (StoreID, EquipmentID) Values (?,?)", s.StoreID, s.EquipmentID)
	if err != nil {
		log.Println(err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	s.StoreItemID = int(id)
	return nil
}

func (s *StoreItem) GetDetails() error {
	err := db.Select(s, "Select * from StoreItem Where StoreitemID = ? LIMIT 1", s.StoreItemID)
	return err
}

func (s *StoreItem) Update() error {
	_, err := db.Exec("Update StoreItem SET StoreID = ? where ID = ?", s.StoreID, s.StoreItemID)
	return err
}

func (s *StoreItem) Delete() error {
	_, err := db.Exec("Delete from StoreItem Where StoreItemID = ?", s.StoreItemID)
	return err
}

func GetStoreItems() ([]StoreItem, error) {
	var ss []StoreItem
	err := db.Get(&ss, "Select * from StoreItem")
	return ss, err
}

func (s *StoreItem) GetFaults() ([]Fault, error) {
	var result []Fault
	err := db.Select(&result, "Select * from Fault Where StoreitemID = ?", s.StoreItemID)
	return result, err
}

func (s *StoreItem) PostFault(f Fault) (Fault, error) {
	err := f.Insert()
	return f, err
}

type Event struct {
	EventID int
	Name    string
	Start   time.Time
	End     time.Time
	Adress  string
}

func (e *Event) Insert() error {
	res, err := db.Exec("Insert Into Event (Name, Start, End, Adress) Values (?, ?, ?, ?)", e.Name, e.Start, e.End, e.Adress)
	if err != nil {
		log.Println(err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	e.EventID = int(id)
	return nil
}

func (e *Event) GetDetails() error {
	err := db.Get(e, "SELECT * from Event Where EventID = ? Limit 1", e.EventID)
	return err
}

func (e *Event) Update() error {
	_, err := db.Exec("Update Event SET Name = ?, Start = ?, End = ?, Adress = ? where ID = ?", e.Name, e.Start, e.End, e.Adress, e.EventID)
	return err
}

func (e *Event) Delete() error {
	_, err := db.Exec("Delete from Event Where ID = ?", e.EventID)
	return err
}

func (e *Event) GetParticipiants() ([]Participiant, error) {
	var pp []Participiant
	err := db.Select(&pp, "Select * from Participiant Where EventID = ?", e.EventID)
	return pp, err
}

func GetEvents() ([]Event, error) {
	var e []Event
	err := db.Select(&e, "Select * from Event")
	return e, err
}

func GetNextEvent() (Event, error) {
	var e Event
	err := db.Get(&e, "Select * from Event where start > ? Order by start ASC Limit 1", time.Now())
	return e, err
}

type Packinglist struct {
	PackinglistID int
	Name          string
	EventID       int
}

func (p *Packinglist) Insert() error {
	res, err := db.Exec("Insert Into Packinglist (Name, EventID) Values (?, ?)", p.Name, p.EventID)
	if err != nil {
		log.Println(err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	p.PackinglistID = int(id)
	return nil
}

func (p *Packinglist) Update() error {
	_, err := db.Exec("Update Packinglist SET name = ?, EventID = ? where ID = ?", p.Name, p.EventID, p.PackinglistID)
	return err
}

func (p *Packinglist) Delete() error {
	_, err := db.Exec("Delete from Packinglist Where PackinglistID = ?", p.PackinglistID)
	return err
}

type PackinglistItem struct {
	PackinglistID int
	StoreitemID   int
}

func (p *PackinglistItem) Insert() error {
	_, err := db.Exec("Insert Into PackinglistItem (PackinglistID, StoreitemID) Values (?, ?)", p.PackinglistID, p.StoreitemID)
	return err
}

func (p *PackinglistItem) Delete() error {
	_, err := db.Exec("Delete from PackinglistItem Where PackinglistID = ?, StoreitemID = ?", p.PackinglistID, p.StoreitemID)
	return err
}

type Participiant struct {
	UserID    int
	EventID   int
	Arrival   time.Time
	Departure time.Time
}

func (p *Participiant) Insert() error {
	_, err := db.Exec("Insert Into Participiant (UserID, EventID, Arrival, Departure) Values (?, ?, ?, ?)", p.UserID, p.EventID)
	return err
}

func (p *Participiant) Update() error {
	_, err := db.Exec("Update Participiant SET UserID = ?, EventID = ?, Arrival = ?, Departure = ? where UserID = ?, EventID = ?", p.UserID, p.EventID, p.Arrival, p.Departure, p.UserID, p.EventID)
	return err
}

func (p *Participiant) Delete() error {
	_, err := db.Exec("Delete from PackinglistItem Where UserID = ?, EventID = ?", p.UserID, p.EventID)
	return err
}

type Wishlist struct {
	WishlistID int
	Name       string
}

func (w *Wishlist) Insert() error {
	res, err := db.Exec("Insert Into Wishlist (Name) Values (?)", w.Name)
	if err != nil {
		log.Println(err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	w.WishlistID = int(id)
	return nil
}

func (w *Wishlist) Update() error {
	_, err := db.Exec("Update Wishlist SET name = ? where ID = ?", w.Name, w.WishlistID)
	return err
}

func (w *Wishlist) Delete() error {
	_, err := db.Exec("Delete from Wishlist Where ID = ?", w.WishlistID)
	return err
}

type Wishlistitem struct {
	WishlistID  int
	EquipmentID int
	Count       int
}

func (p *Wishlistitem) Insert() error {
	_, err := db.Exec("Insert Into Wishlistitem (WishlistID, EquipmentID) Values (?, ?)", p.WishlistID, p.EquipmentID)
	return err
}

func (p *Wishlistitem) Update() error {
	_, err := db.Exec("Update Wishlistitem SET WishlistID = ?, EquipmentID = ? where WishlistID = ?, EquipmentID = ?", p.WishlistID, p.EquipmentID, p.WishlistID, p.EquipmentID)
	return err
}

func (p *Wishlistitem) Delete() error {
	_, err := db.Exec("Delete from Wishlistitem Where WishlistID = ?, EquipmentID = ?", p.WishlistID, p.EquipmentID)
	return err
}

type FaultStatus int

const (
	FaultStatusNew FaultStatus = 1 + iota
	FaultStatusInRepair
	FaultStatusFixed
	FaultStatusUnfixable
)

type Fault struct {
	FaultID     int
	StoreItemID int
	Status      FaultStatus
	Comment     string
}

func (f *Fault) Insert() error {
	res, err := db.Exec("Insert Into Fault (Status, Comment) Values (?, ?)", f.Status, f.Comment)
	if err != nil {
		log.Println(err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	f.FaultID = int(id)
	return nil
}

func GetFaults() ([]Fault, error) {
	var f []Fault
	err := db.Select(&f, "Select * from Fault")
	return f, err
}

func (f *Fault) Update() error {
	_, err := db.Exec("Update Fault SET Status = ?, Comment = ? where ID = ?", f.Status, f.Comment, f.FaultID)
	return err
}

func (f *Fault) Delete() error {
	_, err := db.Exec("Delete from Wishlist Where ID = ?", f.FaultID)
	return err
}

func (f *Fault) GetDetails() error {
	err := db.Get(&f, "Select * from Fault where FaultId = ? Limit 1", f.FaultID)
	return err
}

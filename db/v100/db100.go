package db100

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chaosvermittlung/funkloch-server/global"
	"github.com/jmoiron/sqlx"
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
	ID       int       `json:"id"`
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
	err := db.Get(u, "SELECT * from User Where ID = ? Limit 1", u.ID)
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
	_, err := db.Exec("UPDATE User SET username = ?, password = ?, email = ?, right = ? WHERE id = ?", u.Username, u.Password, u.Email, u.Right, u.ID)
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
	u.ID = int(id)

	return nil
}

func DeleteUser(id int) error {
	_, err := db.Exec("DELETE FROM User Where ID = ?", id)
	return err
}

type Store struct {
	ID      int
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
	s.ID = int(id)
	return nil
}

func (s *Store) GetDetails() error {
	err := db.Get(s, "SELECT * from Store Where ID = ? Limit 1", s.ID)
	return err
}

func (s *Store) GetManager() (User, error) {
	var u User
	u.ID = s.ID
	err := u.GetDetails()
	return u, err
}

func (s *Store) Update() error {
	_, err := db.Exec("Update Store SET name = ?, adress = ?, manager = ? where ID = ?", s.Name, s.Adress, s.Manager, s.ID)
	return err
}

func (s *Store) Delete() error {
	_, err := db.Exec("Delete from Store Where ID = ?", s.ID)
	return err
}

func (s *Store) GetStoreitems() ([]StoreItem, error) {
	var si []StoreItem
	err := db.Select(&si, "Select * from StoreItem Where StoreID = ?", s.ID)
	return si, err
}

type Equipment struct {
	ID   int
	Name string
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
	e.ID = int(id)
	return nil
}

func (e *Equipment) Update() error {
	_, err := db.Exec("Update Equipment SET name = ? where ID = ?", e.Name, e.ID)
	return err
}

func (e *Equipment) Delete() error {
	_, err := db.Exec("Delete from Equipment Where ID = ?", e.ID)
	return err
}

type StoreItem struct {
	ID          int
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
	s.ID = int(id)
	return nil
}

func (s *StoreItem) Update() error {
	_, err := db.Exec("Update StoreItem SET StoreID = ? where ID = ?", s.StoreID, s.ID)
	return err
}

func (s *StoreItem) Delete() error {
	_, err := db.Exec("Delete from Equipment Where ID = ?", s.ID)
	return err
}

type Event struct {
	ID     int
	Name   string
	Start  time.Time
	End    time.Time
	Adress string
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
	e.ID = int(id)
	return nil
}

func (e *Event) Update() error {
	_, err := db.Exec("Update Equipment SET Name = ?, Start = ?, End = ?, Adress = ? where ID = ?", e.Name, e.Start, e.End, e.Adress, e.ID)
	return err
}

func (e *Event) Delete() error {
	_, err := db.Exec("Delete from Event Where ID = ?", e.ID)
	return err
}

func (e *Event) GetParticipiants() ([]Participiant, error) {
	var pp []Participiant
	err := db.Select(&pp, "Select * from Participiant Where EventID = ?", e.ID)
	return pp, err
}

type Packinglist struct {
	ID      int
	Name    string
	EventID int
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
	p.ID = int(id)
	return nil
}

func (p *Packinglist) Update() error {
	_, err := db.Exec("Update Packinglist SET name = ?, EventID = ? where ID = ?", p.Name, p.EventID, p.ID)
	return err
}

func (p *Packinglist) Delete() error {
	_, err := db.Exec("Delete from Packinglist Where ID = ?", p.ID)
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

func (p *PackinglistItem) Update() error {
	_, err := db.Exec("Update PackinglistItem SET PackinglistID = ?, StoreitemID = ? where PackinglistID = ?, StoreitemID = ?", p.PackinglistID, p.StoreitemID, p.PackinglistID, p.StoreitemID)
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
	ID   int
	Name string
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
	w.ID = int(id)
	return nil
}

func (w *Wishlist) Update() error {
	_, err := db.Exec("Update Wishlist SET name = ? where ID = ?", w.Name, w.ID)
	return err
}

func (w *Wishlist) Delete() error {
	_, err := db.Exec("Delete from Wishlist Where ID = ?", w.ID)
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

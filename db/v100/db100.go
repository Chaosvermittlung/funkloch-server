package db100

import (
	"log"
	"strconv"
	"time"

	"github.com/Chaosvermittlung/funkloch-server/global"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/mattn/go-sqlite3"
)

var db *gorm.DB

func Initialisation(dbc *global.DBConnection) {
	var err error
	db, err = gorm.Open(dbc.Driver, dbc.Connection)
	if err != nil {
		log.Fatal(err)
	}
	initDB(dbc)
}

func initDB(dbc *global.DBConnection) {
	var cont bool
	switch dbc.Driver {
	case "sqlite3":
		var err error
		cont, err = global.Exists(dbc.Connection)
		if err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("DB Driver unkown. Stopping Server")
	}
	if !cont {
		db.AutoMigrate(&User{})
		db.AutoMigrate(&Store{})
		db.AutoMigrate(&Equipment{})
		db.AutoMigrate(&Box{})
		db.AutoMigrate(&Item{})
		db.AutoMigrate(&Event{})
		db.AutoMigrate(&Packinglist{})
		db.AutoMigrate(&Participant{})
		db.AutoMigrate(&Wishlist{})
		db.AutoMigrate(&Fault{})

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
		db.Create(&u)
	}

}

type UserRight int

const (
	USERRIGHT_MEMBER UserRight = 1 + iota
	USERRIGHT_ADMIN
)

type User struct {
	gorm.Model
	UserID   int       `json:"id";gorm:"primary_key"`
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
	var u User
	err := db.Where("Username = ?", username).First(&u)
	b := (u.UserID > 0)
	return b, err.Error
}

func GetUsers() ([]User, error) {
	var u []User
	err := db.Find(&u)
	return u, err.Error
}

func (u *User) GetDetailstoUsername() error {
	err := db.Where("Username = ?", u.Username).First(&u)
	return err.Error
}

func (u *User) GetDetails() error {
	err := db.First(&u, u.UserID)
	return err.Error
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
	err := db.Save(&u)
	return err.Error
}

func (u *User) Insert() error {
	err := db.Create(&u)
	return err.Error
}

func DeleteUser(id int) error {
	var u User
	u.UserID = id
	err := db.Delete(&u)
	return err.Error
}

type Store struct {
	gorm.Model
	StoreID   int `gorm:"primary_key"`
	Name      string
	Adress    string `gorm:"foreignkey:ManagerID"`
	Manager   User
	ManagerID int
}

func (s *Store) Insert() error {
	err := db.Create(&s)
	return err.Error
}

func GetStores() ([]Store, error) {
	var s []Store
	err := db.Find(&s)
	return s, err.Error
}

func (s *Store) GetDetails() error {
	err := db.First(&s, s.StoreID)
	return err.Error
}

func (s *Store) GetManager() (User, error) {
	var u User
	err := db.Model(&s).Related(&u, "Manager")
	return u, err.Error
}

func (s *Store) Update() error {
	err := db.Save(&s)
	return err.Error
}

func (s *Store) Delete() error {
	err := db.Delete(&s)
	return err.Error
}

func (s *Store) GetStoreBoxes() ([]Box, error) {
	var bo []Box
	err := db.Model(&s).Related(&bo)
	return bo, err.Error
}

type Equipment struct {
	gorm.Model
	EquipmentID int `gorm:"primary_key"`
	Name        string
}

func (e *Equipment) Insert() error {
	err := db.Create(&e)
	return err.Error
}

func GetEquipment() ([]Equipment, error) {
	var e []Equipment
	err := db.Find(&e)
	return e, err.Error
}

func (e *Equipment) GetDetails() error {
	err := db.First(&e, e.EquipmentID)
	return err.Error
}

func (e *Equipment) Update() error {
	err := db.Save(&e)
	return err.Error
}

func (e *Equipment) Delete() error {
	err := db.Delete(&e)
	return err.Error
}

type Box struct {
	gorm.Model
	BoxID       int `gorm:"primary_key"`
	Store       Store
	Items       []Item
	Code        int
	Description string
}

func (b *Box) Insert() error {
	err := db.Create(&b)
	tmp, err2 := strconv.Atoi(global.CreateBoxEAN(b.BoxID))
	if err2 != nil {
		return err2
	}
	b.Code = tmp
	err = db.Save(&b)
	return err.Error
}

func (b *Box) Update() error {
	err := db.Save(&b)
	return err.Error
}

func GetBoxes() ([]Box, error) {
	var b []Box
	err := db.Find(&b)
	return b, err.Error
}

func (b *Box) GetDetails() error {
	err := db.First(&b, b.BoxID)
	return err.Error
}

func (b *Box) Delete() error {
	err := db.Delete(b)
	return err.Error
}

func (b *Box) AddStoreItem(item Item) error {
	err := db.Model(&b).Association("Items").Append(&item)
	return err.Error
}

type Item struct {
	gorm.Model
	ItemID      int `gorm:"primary_key"`
	BoxID       int
	EquipmentID int
	Equipment   Equipment
	Code        int
	Faults      []Fault
}

func (i Item) Insert() error {
	err := db.Create(&i)
	if err.Error != nil {
		return err.Error
	}
	tmp, err2 := strconv.Atoi(global.CreateItemEAN(i.ItemID))
	if err2 != nil {
		return err2
	}
	i.Code = tmp
	err = db.Save(&i)
	return err.Error
}

func (i Item) GetDetails() error {
	err := db.First(&i, i.ItemID)
	return err.Error
}

func (i Item) Update() error {
	err := db.Save(&i)
	return err.Error
}

func (i Item) Delete() error {
	err := db.Delete(&i)
	return err.Error
}

func GetStoreItems() ([]Item, error) {
	var ii []Item
	err := db.Find(&ii)
	return ii, err.Error
}

func (i Item) GetFaults() ([]Fault, error) {
	var result []Fault
	err := db.Model(&i).Related(&result)
	return result, err.Error
}

func (i Item) PostFault(f Fault) (Fault, error) {
	err := f.Insert()
	return f, err
}

type Event struct {
	gorm.Model
	EventID int `gorm:"primary_key"`
	Name    string
	Start   time.Time
	End     time.Time
	Adress  string
}

func (e *Event) Insert() error {
	err := db.Create(&e)
	return err.Error
}

func (e *Event) GetDetails() error {
	err := db.First(&e, e.EventID)
	return err.Error
}

func (e *Event) Update() error {
	err := db.Save(&e)
	return err.Error
}

func (e *Event) Delete() error {
	err := db.Delete(&e)
	return err.Error
}

func (e *Event) GetParticipants() ([]Participant, error) {
	var pp []Participant
	err := db.Model(&e).Related(&pp)
	return pp, err.Error
}

func (e *Event) GetPackinglists() ([]Packinglist, error) {
	var pp []Packinglist
	err := db.Model(&e).Related(&pp)
	return pp, err.Error
}

func GetEvents() ([]Event, error) {
	var e []Event
	err := db.Find(&e)
	return e, err.Error
}

func GetNextEvent() (Event, error) {
	var e Event
	err := db.Where("start > ?", time.Now()).Order("start asc").First(&e)
	return e, err.Error
}

type Packinglist struct {
	gorm.Model
	PackinglistID int `gorm:"primary_key"`
	Name          string
	EventID       int
	Event         Event
	Boxes         []Box `gorm:"many2many:packinglist_boxes;"`
}

func (p *Packinglist) Insert() error {
	err := db.Create(&p)
	return err.Error
}

func GetPackinglists() ([]Packinglist, error) {
	var p []Packinglist
	err := db.Find(&p)
	return p, err.Error
}

func (p *Packinglist) GetDetails() error {
	err := db.First(&p, p.PackinglistID)
	return err.Error
}

func (p *Packinglist) Update() error {
	err := db.Save(&p)
	return err.Error
}

func (p *Packinglist) GetBoxes() ([]Box, error) {
	var res []Box
	err := db.Model(&p).Related(&res)
	return res, err.Error
}

func (p *Packinglist) Delete() error {
	err := db.Delete(&p)
	return err.Error
}

type Participant struct {
	gorm.Model
	UserID    int `gorm:"primary_key"`
	User      User
	EventID   int `gorm:"primary_key"`
	Event     Event
	Arrival   time.Time
	Departure time.Time
}

func (p *Participant) Insert() error {
	err := db.Create(&p)
	return err.Error
}

func (p *Participant) Update() error {
	err := db.Save(&p)
	return err.Error
}

func (p *Participant) Delete() error {
	err := db.Delete(&p)
	return err.Error
}

func (p *Participant) GetDetails() error {
	err := db.Where("UserID = ? and EventID = ?", p.UserID, p.EventID).First(&p)
	return err.Error
}

type Wishlist struct {
	gorm.Model
	WishlistID int `gorm:"primary_key"`
	Name       string
	Items      []Equipment `gorm:"many2many:wishlist_equipment;"`
}

func (w *Wishlist) Insert() error {
	err := db.Create(&w)
	return err.Error
}

func GetWishlists() ([]Wishlist, error) {
	var ww []Wishlist
	err := db.Find(&ww)
	return ww, err.Error
}

func (w *Wishlist) Update() error {
	err := db.Save(&w)
	return err.Error
}

func (w *Wishlist) Delete() error {
	err := db.Delete(&w)
	return err.Error
}

func (w *Wishlist) GetDetails() error {
	err := db.First(&w, w.WishlistID)
	return err.Error
}

func (w *Wishlist) GetItems() ([]Equipment, error) {
	var res []Equipment
	err := db.Model(&w).Related(&res)
	return res, err.Error
}

type FaultStatus int

const (
	FaultStatusNew FaultStatus = 0 + iota
	FaultStatusInRepair
	FaultStatusFixed
	FaultStatusUnfixable
)

type Fault struct {
	gorm.Model
	FaultID int `gorm:"primary_key"`
	ItemID  int
	Status  FaultStatus
	Comment string
}

func (f *Fault) Insert() error {
	err := db.Create(&f)
	return err.Error
}

func GetFaults() ([]Fault, error) {
	var f []Fault
	err := db.Find(&f)
	return f, err.Error
}

func (f *Fault) Update() error {
	err := db.Save(&f)
	return err.Error
}

func (f *Fault) Delete() error {
	err := db.Delete(&f)
	return err.Error
}

func (f *Fault) GetDetails() error {
	err := db.First(&f, f.FaultID)
	return err.Error
}

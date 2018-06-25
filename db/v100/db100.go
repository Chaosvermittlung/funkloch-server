package db100

import (
	"log"
	"strconv"
	"time"

	"github.com/Chaosvermittlung/funkloch-server/global"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func Initialisation(dbc *global.DBConnection) {
	var err error
	cont := checkDBExists(dbc)
	db, err = gorm.Open(dbc.Driver, dbc.Connection)
	if err != nil {
		log.Fatal(err)
	}
	if !cont {
		initDB()
	}
}

func checkDBExists(dbc *global.DBConnection) bool {
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
	return cont
}

func initDB() {
	//db.LogMode(true)
	log.Println("Creating DB")
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

type UserRight int

const (
	USERRIGHT_MEMBER UserRight = 1 + iota
	USERRIGHT_ADMIN
)

type User struct {
	UserID   int       `json:"id" gorm:"primary_key;AUTO_INCREMENT;not null"`
	Username string    `json:"username" gorm:"not null"`
	Password string    `json:"password" gorm:"not null"`
	Salt     string    `json:"-" gorm:"not null"`
	Email    string    `json:"email" gorm:"not null"`
	Right    UserRight `json:"userright" gorm:"not null"`
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
	if gorm.IsRecordNotFoundError(err.Error) {
		return false, nil
	} else {
		return b, err.Error
	}
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
	StoreID   int    `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Name      string `gorm:"not null"`
	Adress    string `gorm:"not null"`
	Manager   User   `gorm:"not null"`
	ManagerID int    `gorm:"foreignkey:ManagerID;not null"`
	Boxes     []Box  `gorm:"foreignkey:StoreID;association_foreignkey:StoreID"`
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
	err := db.Model(&s).Related(&u)
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
	err := db.Model(&s).Association("Boxes").Find(&bo)
	return bo, err.Error
}

func (s *Store) AddStoreBox(b Box) error {
	err := db.Model(&s).Association("Boxes").Append(&b)
	return err.Error
}

type Equipment struct {
	EquipmentID int    `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Name        string `gorm:"not null"`
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
	BoxID       int    `gorm:"primary_key;AUTO_INCREMENT;not null"`
	StoreID     int    `gorm:"not null"`
	Items       []Item `gorm:"foreignkey:BoxID;association_foreignkey:BoxID"`
	Code        int    `gorm:"type:integer(13)"`
	Description string `gorm:"not null"`
}

type BoxlistEntry struct {
	BoxID            int
	BoxCode          int
	BoxDescription   string
	StoreID          int
	StoreName        string
	StoreAddress     string
	StoreManagerID   int
	StoreManagerName string
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

func (b *Box) AddBoxItem(item Item) error {
	err := db.Model(&b).Association("Items").Append(&item)
	return err.Error
}

func (b *Box) GetBoxItems() ([]Item, error) {
	var ii []Item
	err := db.Model(&b).Association("Items").Find(&ii)
	return ii, err.Error
}

func GetBoxesJoined() ([]BoxlistEntry, error) {
	var ble []BoxlistEntry
	err := db.Table("Boxes").
		Select("Boxes.box_id, Boxes.code as BoxCode, Boxes.description as BoxDescription, Stores.store_id, Stores.name as Storename, Stores.adress as StoreAddress, Stores.manager_id as StoreManagerID, User.Username as StoreManagerName").
		Joins("left join Stores on Boxes.Store_Id = Stores.Store_Id").
		Joins("left join Users on Stores.Manager_id = Users.User_id").
		Scan(&ble)
	return ble, err.Error
}

func (b *Box) GetBoxItemsJoined() ([]ItemslistEntry, error) {
	var ile []ItemslistEntry
	err := db.Table("Items").
		Select("Items.item_id, Items.code as ItemCode, Boxes.box_id, Boxes.code as BoxCode, Boxes.description as BoxDescription, Stores.store_id, Stores.name as Storename, Stores.adress as StoreAddress, Stores.manager_id as StoreManagerID, Equipment.equipment_id, Equipment.name as EquipmentName ").
		Joins("left join Boxes on Items.box_id = Boxes.box_id").
		Joins("left join Stores on Boxes.Store_Id = Stores.Store_Id").
		Joins("left join equipment on Items.equipment_id = equipment.equipment_id").
		Where("Items.Box_id = ?", b.BoxID).
		Scan(&ile)
	return ile, err.Error
}

type Item struct {
	ItemID      int `gorm:"primary_key;AUTO_INCREMENT;not null"`
	BoxID       int
	EquipmentID int       `gorm:"not null"`
	Equipment   Equipment `gorm:"not null"`
	Code        int       `gorm:"type:integer(13)"`
	Faults      []Fault   `gorm:"foreignkey:ItemID;association_foreignkey:ItemID"`
}

type ItemslistEntry struct {
	ItemID         int
	ItemCode       int
	BoxID          int
	BoxCode        int
	BoxDescription string
	StoreID        int
	StoreName      string
	StoreAddress   string
	StoreManagerID int
	EquipmentID    int
	EquipmentName  string
}

func (i *Item) Insert() error {
	err := db.Create(&i)
	tmp, err2 := strconv.Atoi(global.CreateItemEAN(i.ItemID))
	if err2 != nil {
		return err2
	}
	i.Code = tmp
	err = db.Save(&i)
	return err.Error
}

func (i *Item) GetDetails() error {
	err := db.First(&i, i.ItemID)
	return err.Error
}

func (i *Item) GetFullDetails() (ItemslistEntry, error) {
	var ile ItemslistEntry
	err := db.Table("Items").
		Select("Items.item_id, Items.code as ItemCode, Boxes.box_id, Boxes.code as BoxCode, Boxes.description as BoxDescription, Stores.store_id, Stores.name as Storename, Stores.adress as StoreAddress, Stores.manager_id as StoreManagerID, Equipment.equipment_id, Equipment.name as EquipmentName ").
		Joins("left join Boxes on Items.box_id = Boxes.box_id").
		Joins("left join Stores on Boxes.Store_Id = Stores.Store_Id").
		Joins("left join equipment on Items.equipment_id = equipment.equipment_id").
		Where("Items.item_id = ?", i.ItemID).
		Find(&ile)
	return ile, err.Error
}

func (i *Item) Update() error {
	err := db.Save(&i)
	return err.Error
}

func (i *Item) Delete() error {
	err := db.Delete(&i)
	return err.Error
}

func GetItems() ([]Item, error) {
	var ii []Item
	err := db.Find(&ii)
	return ii, err.Error
}

func GetItemsJoined() ([]ItemslistEntry, error) {
	var ile []ItemslistEntry
	err := db.Table("Items").
		Select("Items.item_id, Items.code as ItemCode, Boxes.box_id, Boxes.code as BoxCode, Boxes.description as BoxDescription, Stores.store_id, Stores.name as Storename, Stores.adress as StoreAddress, Stores.manager_id as StoreManagerID, Equipment.equipment_id, Equipment.name as EquipmentName ").
		Joins("left join Boxes on Items.box_id = Boxes.box_id").
		Joins("left join Stores on Boxes.Store_Id = Stores.Store_Id").
		Joins("left join equipment on Items.equipment_id = equipment.equipment_id").
		Scan(&ile)
	return ile, err.Error
}

func (i *Item) GetFaults() ([]Fault, error) {
	var result []Fault
	err := db.Model(&i).Related(&result)
	return result, err.Error
}

func (i *Item) AddFault(f Fault) (Fault, error) {
	err := f.Insert()
	return f, err
}

type Event struct {
	EventID      int           `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Name         string        `gorm:"not null"`
	Start        time.Time     `gorm:"not null"`
	End          time.Time     `gorm:"not null"`
	Adress       string        `gorm:"not null"`
	Participants []Participant `gorm:"foreignkey:EventID;association_foreignkey:EventID"`
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
	PackinglistID int    `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Name          string `gorm:"not null"`
	EventID       int    `gorm:"not null"`
	Event         Event  `gorm:"not null"`
	Boxes         []Box  `gorm:"many2many:packinglist_boxes;"`
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

func (p *Packinglist) AddPackinglistBox(b Box) error {
	err := db.Model(&p).Association("Boxes").Append(&b)
	return err.Error
}

func (p *Packinglist) GetPackinglistBoxes() ([]Box, error) {
	var res []Box
	err := db.Model(&p).Association("Boxes").Find(&res)
	return res, err.Error
}

func (p *Packinglist) RemovePackinglistBox(b Box) error {
	err := db.Model(&p).Association("Boxes").Delete(&b)
	return err.Error
}

func (p *Packinglist) Delete() error {
	err := db.Delete(&p)
	return err.Error
}

type Participant struct {
	UserID    int       `gorm:"type:integer;primary_key;not null"`
	User      User      `gorm:"not null;foreignkey:UserID;association_foreignkey:UserID"`
	EventID   int       `gorm:"type:integer;primary_key;not null"`
	Arrival   time.Time `gorm:"not null"`
	Departure time.Time `gorm:"not null"`
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
	err := db.Where("User_ID = ? and Event_ID = ?", p.UserID, p.EventID).First(&p)
	return err.Error
}

type Wishlist struct {
	WishlistID int         `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Name       string      `gorm:"not null"`
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

func (w *Wishlist) AddWishlistItem(e Equipment) error {
	err := db.Model(&w).Association("Items").Append(&e)
	return err.Error
}

func (w *Wishlist) GetWishlistItems() ([]Equipment, error) {
	var res []Equipment
	err := db.Model(&w).Association("Items").Find(&res)
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
	FaultID int         `gorm:"primary_key;AUTO_INCREMENT;not null"`
	ItemID  int         `gorm:"not null"`
	Status  FaultStatus `gorm:"not null"`
	Comment string      `gorm:"not null"`
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

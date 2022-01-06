package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName   string `gorm:"size:100;not null" json:"first_name"`
	LastName    string `gorm:"size:100;not null" json:"last_name"`
	Email       string `gorm:"size:100;not null;unique" json:"email"`
	Password    string `gorm:"size:100;not null" json:"password"`
	PhoneNumber string `json:"phone_number"`
	Role        string `gorm:"default:''" json:"role"`
}

type Coach struct {
	ID       int
	UserID   int
	User     *User
	Sessions []Session `gorm:"foreignKey:CoachID"`
	Tasks    []Task    `gorm:"foreignKey:AssignerID"`
}

type Client struct {
	ID       int
	UserID   uint
	User     *User
	Sessions []Session `gorm:"foreignKey:ClientID"`
	Tasks    []Task    `gorm:"foreignKey:AssigneeID"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if pw, err := bcrypt.GenerateFromPassword([]byte(u.Password), 5); err == nil {
		tx.Statement.SetColumn("Password", string(pw))
	}
	return
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (u *User) GetUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	if err = db.Model(&User{}).Find(&users).Error; err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) GetUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	if err = db.Model(User{}).Where("id = ?", uid).First(&u).Error; err != nil {
		return &User{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) CreateUser(db *gorm.DB) (*User, error) {
	var err error
	var c Coach
	var cl Client
	if err = db.Debug().Create(&u).Error; err != nil {
		return &User{}, err
	}
	switch u.Role {
	case "coach":
		c.User = u
		c.ID = int(u.ID)
		c.UserID = int(u.ID)
		db.Debug().Create(&c)
	case "client":
		cl.User = u
		cl.ID = int(u.ID)
		db.Debug().Create(&cl)
	default:
		return u, nil
	}
	return u, nil
}

func (u *User) DeleteUser(db *gorm.DB, uid uint32) (int64, error) {
	if error := db.Model(&User{}).Where("id = ?", uid).Delete(&User{}).Error; error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *User) GetClients(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	if err = db.Model(&User{}).Find(&users, "role = ?", "client").Error; err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) GetClientByID(db *gorm.DB, uid int) (*User, error) {
	var err error
	if err = db.Model(&User{}).Where("id = ? AND role = ?", uid, "client").First(&u).Error; err != nil {
		return &User{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, errors.New("Client Not Found")
	}
	return u, err
}

func (u *User) GetCoaches(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	if err = db.Model(&User{}).Find(&users, "role = ?", "coach").Error; err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) GetCoachByID(db *gorm.DB, uid int) (*User, error) {
	var err error
	if err = db.Model(&User{}).Where("id = ? AND role = ?", uid, "coach").First(&u).Error; err != nil {
		return &User{}, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, errors.New("Coach Not Found")
	}
	return u, err
}

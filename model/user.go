package model

import (
	"dairy_service/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html"
	"strings"
)

// The User struct is composed of the Gorm Model struct, a string for the user’s username, another string for the password,
// and a slice of Entry items. In this way, you're specifying a one-to-many relationship between the User struct and the Entry structs.
type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"-"`
	Entries  []Entry
}

// The Save function adds a new user to the database (in the absence of any errors). Before saving,
// any whitespace in the provided username is trimmed out and the provided password is hashed for security purposes.
func (user *User) Save() (*User, error) {
	err := database.Database.Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (user *User) BeforeSave(*gorm.DB) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passwordHash)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	return nil
}

func (user *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func FindUserByUsername(username string) (User, error) {
	var user User
	err := database.Database.Where("username=?", username).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// FindUserById - In addition to retrieving the user, the entries associated with the user are eagerly loaded - t
// hus populating the Entries slice in the User struct.
func FindUserById(id uint) (User, error) {
	var user User
	err := database.Database.Preload("Entries").Where("ID=?", id).Find(&user).Error
	if err != nil {
		return User{}, err
	}
	return user, nil
}

package authorization

import (
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"longpollServer/utils"
)

type Controller struct {
	db        *gorm.DB
	tokenAuth *jwtauth.JWTAuth
}

func NewController(db *gorm.DB, tokenAuth *jwtauth.JWTAuth) *Controller {
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal(err)
	}
	return &Controller{db: db, tokenAuth: tokenAuth}
}

func (c *Controller) Login(login, password string) (string, string, error) {
	var user User
	err := c.db.Where("email = ?", login).First(&user).Error
	if err != nil {
		return "", "", err
	}

	if user.Password != password {
		return "", "", errors.New("invalid login/password")
	}

	_, tokenString, err := c.tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})
	if err != nil {
		return "", "", err
	}

	return user.ID.String(), tokenString, nil
}

func (c *Controller) CreateUser(username, email, password string) (string, string, error) {
	if !utils.IsValidEmail(email) || !utils.IsValidPassword(password) {
		return "", "", errors.New("invalid email/password")
	}

	err := c.db.Where("email = ?", email).First(&User{}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", "", errors.New("user already exists")
	}

	user := User{
		ID:       uuid.New(),
		Username: username,
		Email:    email,
		Password: password,
	}
	err = c.db.Create(&user).Error
	if err != nil {
		return "", "", err
	}
	_, tokenString, err := c.tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})
	if err != nil {
		return "", "", err
	}

	return user.ID.String(), tokenString, nil
}

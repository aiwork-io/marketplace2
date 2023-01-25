package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const CTX_KEY_USER string = "models.user"

type User struct {
	Id         string     `json:"id" gorm:"primaryKey;size:66"`
	Name       string     `json:"name" gorm:"size:256"`
	Email      string     `json:"email" gorm:"unique;size:256"`
	Wallet     string     `json:"wallet" gorm:"size:128"`
	Password   string     `json:"password" gorm:"size:256"`
	VerifiedAt *time.Time `json:"verified_at"`
	CreatedAt  *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  *time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Role int `json:"role" gorm:"default:0"`
}

func (u *User) IsAdmin() bool {
	return u.Role == 1
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) SignToken(secret string) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     u.Id,
		"name":   u.Name,
		"email":  u.Email,
		"wallet": u.Wallet,
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(secret))
}

func ValidateToken(secret, tokenString string) (*User, error) {
	// Create a new token object, specifying signing method and the claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		u := User{
			Id:     claims["id"].(string),
			Name:   claims["name"].(string),
			Email:  claims["email"].(string),
			Wallet: claims["wallet"].(string),
		}
		return &u, nil
	}

	return nil, errors.New("invalid token")
}

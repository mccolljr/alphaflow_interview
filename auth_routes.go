package main

import (
	"alphaflow/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/attache/attache"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var _TokenSecret = []byte(os.Getenv("TOKEN_SECRET"))

func (c *AlphaFlow) POST_Signup() {
	var body struct {
		Email          string `json:"email"`
		Password       string `json:"password"`
		PasswordVerify string `json:"password_verify"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		log.Println(err)
		attache.ErrorMessageJSON(http.StatusBadRequest, "unable to parse request body")
	}
	if body.Password != body.PasswordVerify {
		attache.ErrorMessageJSON(http.StatusBadRequest, "passwords do not match")
	}
	if body.Password == "" {
		// TODO: validate password strength, etc
		attache.ErrorMessageJSON(http.StatusBadRequest, "password cannot be empty")
	}
	if !strings.Contains(body.Email, "@") {
		// TODO: validate email better
		attache.ErrorMessageJSON(http.StatusBadRequest, "invalid email address")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	if err != nil {
		attache.ErrorFatal(err)
	}

	var u models.User
	u.Email = strings.ToLower(body.Email)
	u.HashPW = string(hashed)
	if err := c.DB().Insert(&u); err != nil {
		attache.ErrorFatal(err)
	}
	attache.RenderJSON(c.ResponseWriter(), u)
}

func (c *AlphaFlow) POST_Login() {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		log.Println(err)
		attache.ErrorMessageJSON(http.StatusBadRequest, "unable to parse request body")
	}
	var user models.User
	if err := c.DB().GetBy(&user, "email = ?", body.Email); err != nil {
		// Don't expose the fact that the user exists (or not)
		attache.ErrorMessageJSON(http.StatusBadRequest, "invalid user/pass combination")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPW), []byte(body.Password)); err != nil {
		log.Println("login: bad password for", body.Email)
		attache.ErrorMessageJSON(http.StatusBadRequest, "invalid user/pass combination")
	}

	tok, err := c.getToken(user)
	if err != nil {
		attache.ErrorFatal(err)
	}

	attache.RenderJSON(c.ResponseWriter(), map[string]interface{}{"token": tok, "user": user})
}

func (c *AlphaFlow) getToken(user models.User) (string, error) {
	// TODO: strengthen the signing method
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
		"id":  user.ID,
	})
	return tok.SignedString(_TokenSecret)
}

func (c *AlphaFlow) getUser() (*models.User, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, fmt.Errorf("invalid auth header")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	token, parseErr := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return _TokenSecret, nil
	})
	if parseErr != nil {
		return nil, fmt.Errorf("parse token: %w", parseErr)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("validate token: %w", parseErr)
	}

	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		return nil, fmt.Errorf("expired token")
	}

	if claims["id"] == nil {
		return nil, fmt.Errorf("malformed claims")
	}

	var target models.User
	if err := c.DB().Get(&target, claims["id"]); err != nil {
		if err == attache.ErrRecordNotFound {
			return nil, fmt.Errorf("no such user")
		}
		return nil, fmt.Errorf("load user: %w", err)
	}

	return &target, nil
}

package main

import (
	"alphaflow/models"
	"log"
	"net/http"
	"strings"

	"github.com/attache/attache"
	"github.com/dgrijalva/jwt-go"
)

const _AuthHeaderPrefix = "Bearer "

func (c *AlphaFlow) loadUserFromToken() {
	authHeader := c.Request().Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, _AuthHeaderPrefix) {
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, _AuthHeaderPrefix)
	token, parseErr := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret-key"), nil
	})
	if parseErr != nil {
		log.Println("parse token:", parseErr)
		return
	}
	if validErr := token.Claims.Valid(); validErr != nil {
		log.Println("validate token:", validErr)
		return
	}
	id, ok := token.Claims.(jwt.MapClaims)["id"].(string)
	if !ok {
		log.Println("invalid token claims")
		return
	}
	var target models.User
	if err := c.DB().Get(&target, id); err != nil {
		if err == attache.ErrRecordNotFound {
			log.Println("invalid token user")
			return
		}
		log.Println("error loading token user:", err)
		return
	}
	c.User = &target
}

func (c *AlphaFlow) requireUser() {
	if c.User == nil {
		attache.ErrorMessage(http.StatusUnauthorized, "invalid auth")
	}
}

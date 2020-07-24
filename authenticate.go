package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// AccountAuthenticate authenticates account provided via HTTP POST request using data from DB.
func AccountAuthenticate(c *gin.Context) {
	var inputAccount, dbAccount Account

	// decode input data to object
	if err := c.ShouldBind(&inputAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	// decode DB data to object
	if err := db.Collection(collectionName).FindOne(
		context.TODO(),
		bson.M{
			"email": inputAccount.Email,
		},
	).Decode(&dbAccount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	// validate password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(dbAccount.Password),
		[]byte(inputAccount.Password),
	); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "authorized"})
}

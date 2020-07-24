package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

// AccountCreate creates account from data provided via HTTP POST request.
func AccountCreate(c *gin.Context) {
	var inputAccount, dbAccount Account

	// decode input data to object
	if err := c.ShouldBind(&inputAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	// decode DB data to object
	_ = db.Collection(collectionName).FindOne(
		context.TODO(),
		bson.M{
			"email": inputAccount.Email,
		},
	).Decode(&dbAccount)

	// check if account already exist
	if strings.EqualFold(inputAccount.Email, dbAccount.Email) {
		c.JSON(http.StatusForbidden, gin.H{"status": "account exist"})

		return
	}

	// get salted account data
	dbAccount = inputAccount.GetSecureAccount()

	// encode object to BSON
	b, err := bson.Marshal(dbAccount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	// insert data to collenction
	_, err = db.Collection(collectionName).InsertOne(context.Background(), b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

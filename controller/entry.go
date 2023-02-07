package controller

import (
	"dairy_service/helper"
	"dairy_service/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// The AddEntry function marshals the body of the request into an Entry struct after which it gets the
// currently authenticated user from the request header. Next, it sets the associated user ID for the entry and saves it.
// The saved details are then returned in a JSON response.
func AddEntry(context *gin.Context) {
	var input model.Entry
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := helper.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserID = user.ID

	savedEntry, err := input.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"data": savedEntry})
}

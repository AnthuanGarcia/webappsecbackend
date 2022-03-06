package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"

	db "github.com/AnthuanGarcia/appWebSeguridad/db"
	helper "github.com/AnthuanGarcia/appWebSeguridad/src/helpers"
	models "github.com/AnthuanGarcia/appWebSeguridad/src/models"
)

var validate = validator.New()

//CreateUser is the api used to tget a single user
func SignUp(c *gin.Context) {

	user := &models.User{}

	if err := c.BindJSON(user); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	validationErr := validate.Struct(user)
	if validationErr != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return

	}

	if !helper.ValidatePassword(user.Contraseña) {

		c.JSON(http.StatusBadRequest, gin.H{"error": "La contraseña debe contener una letra minuscula, una letra mayuscula, un numero y un simbolo"})
		return

	}

	user.ID = primitive.NewObjectID()

	user.Contraseña = helper.HashPassword(user.Contraseña)

	user.Fch_Creacion, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Fch_Renovacion, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	token, refreshToken, err := helper.GenerateAllTokens(
		user.Username,
		user.Nombre,
		user.Ape_Paterno,
		user.ID.Hex(),
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	user.Token = token
	user.Token_Act = refreshToken

	result, err := db.InsertUser(user)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, result)

}

//Login is the api used to tget a single user
func Login(c *gin.Context) {

	user := &models.User{}

	if err := c.BindJSON(user); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	verifyUser, err := db.GetUser(user)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	token, refreshToken, err := helper.GenerateAllTokens(
		verifyUser.Username,
		verifyUser.Nombre,
		verifyUser.Ape_Paterno,
		verifyUser.ID.Hex(),
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	err = db.UpdateAllTokens(token, refreshToken, verifyUser.ID)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, verifyUser)

}

func Dashboard(c *gin.Context) {

	users, err := db.GetUsers()

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, users)

}

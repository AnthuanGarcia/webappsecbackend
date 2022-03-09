package controllers

import (
	"encoding/json"
	"io/ioutil"
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

	c.JSON(http.StatusOK, verifyUser.Token)

}

func Dashboard(c *gin.Context) {

	users, err := db.GetUsers()

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, users)

}

func ManyUsers(c *gin.Context) {

	users := []models.User{}

	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	err = json.Unmarshal(body, &users)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	for i := range users {

		validationErr := validate.Struct(users[i])
		if validationErr != nil {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": validationErr.Error(),
				"usr":   users[i],
			})

			return

		}

		if !helper.ValidatePassword(users[i].Contraseña) {

			c.JSON(http.StatusBadRequest, gin.H{"error": "La contraseña debe contener una letra minuscula, una letra mayuscula, un numero y un simbolo"})
			return

		}

		users[i].ID = primitive.NewObjectID()

		users[i].Contraseña = helper.HashPassword(users[i].Contraseña)

		users[i].Fch_Creacion, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		users[i].Fch_Renovacion, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		token, refreshToken, err := helper.GenerateAllTokens(
			users[i].Username,
			users[i].Nombre,
			users[i].Ape_Paterno,
			users[i].ID.Hex(),
		)

		if err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}

		users[i].Token = token
		users[i].Token_Act = refreshToken

	}

	userModified := make([]interface{}, len(users))

	for i, user := range users {
		userModified[i] = user
	}

	err = db.InsertManyUsers(userModified)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, gin.H{"ok": "Usuarios insertados"})

}

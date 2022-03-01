package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	db "github.com/AnthuanGarcia/appWebSeguridad/db"
	helper "github.com/AnthuanGarcia/appWebSeguridad/src/helpers"
	models "github.com/AnthuanGarcia/appWebSeguridad/src/models"
)

var (
	//userCollection *mongo.Collection = db.GetCollection(db.COLLECTION)
	validate = validator.New()
)

// hashPassword -  used to encrypt the password before it is stored in the DB
func hashPassword(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(bytes)

}

/*VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg

}*/

//CreateUser is the api used to tget a single user
func SignUp(c *gin.Context) {

	//var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
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

	/*count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	defer cancel()
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		return
	}*/

	/*count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	defer cancel()
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
		return
	}*/

	/*resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		msg := fmt.Sprintf("User item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	defer cancel()
	*/

	user.ID = primitive.NewObjectID()

	user.Contraseña = hashPassword(user.Contraseña)

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

	//var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	user := &models.User{}
	//var foundUser models.User

	if err := c.BindJSON(user); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	/*err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login or passowrd is incorrect"})
		return
	}

	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()
	if passwordIsValid != true {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}*/

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

	err = db.UpdateAllTokens(token, refreshToken, verifyUser.ID.Hex())

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, verifyUser)

}

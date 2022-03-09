package db

import (
	"context"
	"errors"
	_ "fmt"
	"log"
	_ "os"
	"time"

	model "github.com/AnthuanGarcia/appWebSeguridad/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Variables para la conexion
const (
	_TIMEOUT      = 15
	_STR_TEMPLATE = "mongodb+srv://%s:%s@webappsec.odee2.mongodb.net/%s?retryWrites=true&w=majority"
	_URI          = "mongodb://localhost:27017"
	_DB_NAME      = "WebAppSec"
	_COLLECTION   = "Usuarios"
)

// getConnection - Conexion a MongoDB
func connect() (*mongo.Client, context.Context, context.CancelFunc) {

	/*username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	clusterEndPoint := os.Getenv("MONGO_ENDPOINT")

	// Generamos una URI para conectarnos al cliente de Mongo
	connectionURI := fmt.Sprintf(
		_STR_TEMPLATE,
		username,
		password,
		clusterEndPoint,
	)*/

	// Obtenemos un elemento del tipo context para poder realizar consultas
	ctx, cancel := context.WithTimeout(
		context.Background(),
		_TIMEOUT*time.Second,
	)

	// Conectamos al cliente de Mongo con la URI generada
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(_URI))

	if err != nil {
		log.Panicf("Fallo al conectar al cluster %v\n", err)
	}

	//log.Printf("Conectado a MongoDB\n")

	// Retornamos un puntero al cliente, el contexto, y una funcion para cancelar la conexion
	return client, ctx, cancel

}

func InsertUser(user *model.User) (*mongo.InsertOneResult, error) {

	client, ctx, _ := connect()

	defer client.Disconnect(ctx)

	database := client.Database(_DB_NAME)
	collection := database.Collection(_COLLECTION)

	count, err := collection.CountDocuments(ctx, bson.M{"username": user.Username})

	if err != nil {
		return nil, err
	}

	if count >= 1 {
		return nil, errors.New("El nombre de usuario ya esta registrado")
	}

	result, err := collection.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	return result, nil

}

func InsertManyUsers(users []interface{}) error {

	client, ctx, _ := connect()

	defer client.Disconnect(ctx)

	database := client.Database(_DB_NAME)
	collection := database.Collection(_COLLECTION)

	_, err := collection.InsertMany(ctx, users)

	if err != nil {
		return err
	}

	return nil

}

func GetUsers() (users []*model.User, err error) {

	client, ctx, _ := connect()

	defer client.Disconnect(ctx)

	database := client.Database(_DB_NAME)
	collection := database.Collection(_COLLECTION)

	cursor, err := collection.Find(ctx, bson.D{})

	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil

}

func GetUser(user *model.User) (*model.User, error) {

	//TODO: Funcion para obtener documento de un usuario

	client, ctx, _ := connect()

	defer client.Disconnect(ctx)

	database := client.Database(_DB_NAME)
	collection := database.Collection(_COLLECTION)

	result := &model.User{}

	err := collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(result)

	if err != nil {
		return nil, errors.New("Usuario no registrado")
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Contraseña), []byte(user.Contraseña))

	if err != nil {
		return nil, errors.New("Contraseña incorrecta")
	}

	return result, nil

}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId primitive.ObjectID) error {

	client, ctx, _ := connect()

	defer client.Disconnect(ctx)

	database := client.Database(_DB_NAME)
	collection := database.Collection(_COLLECTION)

	updateObj := primitive.D{}

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "token_act", Value: signedRefreshToken})

	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "fch_renovacion", Value: updated_at})

	upsert := true

	filter := bson.M{"_id": userId}

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		&opt,
	)

	if err != nil {
		return err
	}

	return nil

}

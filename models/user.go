package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User - Informacion de usuarios
type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Username      string             `json:"username"    validate:"required,min=2,max=100"`
	Nombre        string             `json:"nombre"      validate:"required,min=2,max=100"`
	ApePaterno    string             `json:"ape_paterno" validate:"required,min=2,max=100"`
	ApeMaterno    string             `json:"ape_materno" validate:"required,min=2,max=100"`
	Telefono      string             `json:"telefono"    validate:"required,min=6,max=9"`
	Direccion     string             `json:"direccion"   validate:"required,min=2,max=100"`
	Contraseña    string             `json:"contraseña"  validate:"required,min=8,max=32"`
	Token         string             `json:"token"`
	TokenAct      string             `json:"token_act"`
	FchCreacion   time.Time          `json:"fch_creacion"`
	FchRenovacion time.Time          `json:"fch_renovacion"`
}

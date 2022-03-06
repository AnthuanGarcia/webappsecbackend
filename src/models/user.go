package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User - Informacion de usuarios
type User struct {
	ID             primitive.ObjectID `bson:"_id" json:"ID"`
	Username       string             `json:"Username"    validate:"required,min=2,max=100"`
	Nombre         string             `json:"Nombre"      validate:"required,min=2,max=100"`
	Ape_Paterno    string             `json:"Ape_paterno" validate:"required,min=2,max=100"`
	Ape_Materno    string             `json:"Ape_materno" validate:"required,min=2,max=100"`
	Telefono       string             `json:"Telefono"    validate:"required,min=6,max=12"`
	Direccion      string             `json:"Direccion"   validate:"required,min=2,max=100"`
	Contraseña     string             `json:"Contraseña"  validate:"required,min=8,max=32"`
	Token          string             `json:"Token"`
	Token_Act      string             `json:"Token_act"`
	Fch_Creacion   time.Time          `json:"Fch_creacion"`
	Fch_Renovacion time.Time          `json:"Fch_renovacion"`
}

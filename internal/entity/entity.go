package entity

import "time"

// админские права 2 юзерские 1
const AdminPermission = 2
const UserPermission = 1

type User struct {
	Id         int
	Login      string
	Password   string
	Permission int
}

const SexMale = "male"
const SexFemale = "female"

// struct of actor
//
// name, data of birthday, sex
//
// swagger:model actor
type Actor struct {
	// id actor int database or in query (different)
	//
	// required: false
	Id int `json:"id,omitempty"`

	// Name of actor
	//
	// required: true
	// example: jon
	Name string `json:"name"`

	// Sex of actor (male or female only)
	//
	// required: true
	// example:male
	Sex string `json:"sex"`

	// data of birthday of actor
	// example: 01-01-1999
	//
	// required: true
	DataOfBirthday time.Time `json:"dataofbirthday"`

	// List films of actor
	//
	// required: false
	// swagger:allOf
	Films []Film `json:"films,omitempty"`
}

// struct of films
//
// film with name, description and data of realease
//
// swagger:model film
type Film struct {

	// id film int database or in query (different)
	//
	// required: false
	// example: 1
	Id int `json:"id,omitempty"`

	// Name of film
	//
	// required: true
	// example: alive
	Name string `json:"name"`

	// Information about film
	//
	// required: true
	// min length: 1
	// min length: 1000
	// example: good film
	About string `json:"about"`

	// Information realease data
	//
	// required: true
	// example: 02-01-1999
	ReleaseDate time.Time `json:"releasedate"`

	// Raiting of actor
	//
	// required: true
	// Minimum 0
	// Maximum 10
	// example: 4
	Rating int `json:"rating"`

	// List of actor in this film
	//
	// required: false
	// swagger:allOf
	Actors []Actor `json:"actors,omitempty"`
}

// swagger:model editactor
type RequestEditActor struct {
	// swagger:allOf
	Oldactor Actor `json:"oldactor"`
	// swagger:allOf
	Newactor Actor `json:"newactor"`
}

// swagger:model editfilm
type RequestEditFilm struct {
	// swagger:allOf
	Oldfilm Film `json:"oldfilm"`
	// swagger:allOf
	NewFilm Film `json:"newfilm"`
}

// swagger:model connection
type RequestEditConnection struct {
	// swagger:allOf
	Film Film `json:"film"`
	// swagger:allOf
	Actor Actor `json:"actor"`
}

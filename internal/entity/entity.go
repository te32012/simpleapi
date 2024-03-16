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

type Actor struct {
	Id             int       `json:"id,omitempty"`
	Name           string    `json:"name"`
	Sex            string    `json:"sex"`
	DataOfBirthday time.Time `json:"dataofbirthday"`
	Films          []Film    `json:"films,omitempty"`
}

type Film struct {
	Id          int       `json:"id,omitempty"`
	Name        string    `json:"name"`
	About       string    `json:"about"`
	ReleaseDate time.Time `json:"releasedate"`
	Rating      int       `json:"rating"`
	Actors      []Actor   `json:"actors,omitempty"`
}

type RequestEditActor struct {
	Oldactor Actor `json:"oldactor"`
	Newactor Actor `json:"newactor"`
}

type RequestEditFilm struct {
	Oldfilm Film `json:"oldfilm"`
	NewFilm Film `json:"newfilm"`
}

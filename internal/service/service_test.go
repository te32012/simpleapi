package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
	"vktestgo2024/internal/entity"
	"vktestgo2024/internal/service"
)

var list int = 0
var list2 int = 0

type authmock struct {
}

func (a *authmock) LoginUser(login string, password string) (string, error) {
	return "", nil
}
func (a *authmock) CheckUserIsLoginedAndHasPermission(key string, operation int) (*entity.User, bool) {
	return nil, false
}

type Databasemock struct {
}

func (db *Databasemock) GetUser(ctx context.Context, login string, password string) (*entity.User, error) {
	if "user" == login && "pswd" == password {
		return &entity.User{Id: 1, Login: "user", Password: "pswrd", Permission: entity.UserPermission}, nil
	} else {
		return nil, errors.New("тестовая ошибка")
	}
}
func (db *Databasemock) AddActor(ctx context.Context, actor entity.Actor) error {
	if actor.Name == "jon" {
		return nil
	} else {
		return errors.New("тестовая ошибка add")
	}
}
func (db *Databasemock) EditActor(ctx context.Context, actor entity.Actor) error {
	if actor.Name == "jane" {
		return nil
	} else {
		return errors.New("тестовая ошибка edit")
	}
}
func (db *Databasemock) GetActor(ctx context.Context, actor entity.Actor) (*entity.Actor, error) {
	if actor.Name == "jon" {
		return &actor, nil
	} else {
		return nil, errors.New("тестовая ошибка get")
	}
}
func (db *Databasemock) DeleteActor(ctx context.Context, actor entity.Actor) error {
	if actor.Name == "jon" {
		return nil
	} else {
		return errors.New("тестовая ошибка delete")
	}

}
func (db *Databasemock) GetActors(ctx context.Context) ([]entity.Actor, error) {
	switch list {
	case 1:
		return make([]entity.Actor, 0), nil
	case 2:
		var lst []entity.Actor
		lst = append(lst, entity.Actor{Id: 1, Name: "jane", Sex: entity.SexFemale, DataOfBirthday: time.Time{}})
		return lst, nil
	case 3:
		var lst []entity.Actor
		lst = append(lst, entity.Actor{Id: 2, Name: "jon", Sex: entity.SexFemale, DataOfBirthday: time.Time{}})
		return lst, nil
	default:
		return nil, errors.New("тестовая ошибка getactors")
	}

}
func (db *Databasemock) AddFilm(ctx context.Context, film entity.Film) error {
	if film.Name == "alive" {
		return nil
	} else {
		return errors.New("тестовая ошибка")
	}
}
func (db *Databasemock) GetFilm(ctx context.Context, film entity.Film) (*entity.Film, error) {
	if film.Name == "Warhorse One" {
		return &film, nil
	} else {
		return nil, errors.New("тестовая ошибка getfilm")
	}
}
func (db *Databasemock) DeleteFilm(ctx context.Context, film entity.Film) error {
	if film.Name == "Warhorse One" {
		return nil
	} else {
		return errors.New("тестовая ошибка deletefilm")
	}
}
func (db *Databasemock) EditFilm(ctx context.Context, film entity.Film) error {
	if film.Name == "Warhorse One" {
		return errors.New("тестовая ошибка editfilm")
	} else {
		return nil
	}
}
func (db *Databasemock) GetListFilms(ctx context.Context, keySort int, orderSort int) ([]entity.Film, error) {
	switch list2 {
	case 1:
		return make([]entity.Film, 0), nil
	case 2:
		var films []entity.Film
		films = append(films, entity.Film{Name: "alive", ReleaseDate: time.Time{}, Rating: 10, About: "good film"})
		return films, nil
	default:
		return nil, errors.New("тестовая ошибка getlistfilm")
	}

}
func (db *Databasemock) FindInFilm(ctx context.Context, fragment string) ([]entity.Film, error) {
	if fragment == "Warhorse One" {
		var films []entity.Film
		films = append(films, entity.Film{Id: 1, Name: "Warhorse One", ReleaseDate: time.Time{}, Rating: 10, About: "good film"})
		return films, nil
	}
	if fragment == "alive" {
		var films []entity.Film
		films = append(films, entity.Film{Id: 2, Name: "alive", ReleaseDate: time.Time{}, Rating: 10, About: "good film"})
		return films, nil
	}
	if fragment == "test" {
		var films []entity.Film
		return films, nil
	}
	return nil, errors.New("тестовая ошибка findinfilm")
}
func (db *Databasemock) GetListFilmByActorId(ctx context.Context, id_actor int) ([]entity.Film, error) {
	switch list {
	case 2:
		if id_actor != 1 {
			return nil, errors.New("тестовая ошибка getlistfilmbyactorid")
		} else {
			return nil, nil
		}
	case 3:
		if id_actor != 1 {
			return nil, errors.New("тестовая ошибка getlistfilmbyactorid")
		} else {
			return nil, nil
		}
	}
	return nil, nil
}
func (db *Databasemock) AddActorFilmConnection(ctx context.Context, id_actor, id_film int) error {
	if id_actor == 1 {
		return errors.New("тестовая ошибка")
	}
	return nil
}
func (db *Databasemock) DeleteActorFilmConnection(ctx context.Context, id_actor, id_film int) error {
	if id_actor == 1 {
		return errors.New("тестовая ошибка")
	}
	return nil
}

func TestService(t *testing.T) {
	service := service.NewService(&authmock{}, &Databasemock{})
	s, _ := service.Login("a", "a")
	if s != "" {
		t.Fatal()
	}
	ok := service.CheckUserIsLoginedAndHasPermission("a", 1)
	if ok {
		t.Fatal()
	}

	// тест 1
	var actor entity.Actor
	actor.Id = 0
	actor.Name = "jon"
	actor.Sex = entity.SexMale
	actor.DataOfBirthday = time.Time{}
	data, _ := json.Marshal(actor)
	e := service.AddActor(data)
	if e != nil {
		t.Fatal(e)
	}
	actor.Sex = "s"
	data, _ = json.Marshal(actor)
	e = service.AddActor(data)
	if e == nil {
		t.Fatal(e)
	}
	actor.Sex = entity.SexFemale
	actor.Name = "jane"
	data, _ = json.Marshal(actor)
	e = service.AddActor(data)
	if e == nil {
		t.Fatal(e)
	}
	actor.Sex = entity.SexMale
	actor.Name = "jon"
	data, _ = json.Marshal(actor)
	data = append(data, 1)
	e = service.AddActor(data)
	if e == nil {
		t.Fatal(e)
	}

	// тест 2
	var actor2 entity.Actor

	actor2.Id = 0
	actor2.Name = "jane"
	actor2.Sex = entity.SexFemale
	actor2.DataOfBirthday = time.Now()
	var request entity.RequestEditActor
	request.Oldactor = actor
	request.Newactor = actor2

	data, _ = json.Marshal(request)
	e = service.EditActor(data)
	if e != nil {
		t.Fatal(e)
	}

	request.Newactor = actor

	data, _ = json.Marshal(request)
	e = service.EditActor(data)
	if e == nil {
		t.Fatal(e)
	}
	request.Newactor = actor2

	data, _ = json.Marshal(request)
	data = append(data, 1)

	e = service.EditActor(data)
	if e == nil {
		t.Fatal(e)
	}

	request.Oldactor = actor2

	data, _ = json.Marshal(request)

	e = service.EditActor(data)
	if e == nil {
		t.Fatal(e)
	}

	// тест 3
	actor.Name = "jon"
	data, _ = json.Marshal(actor)

	e = service.DeleteActor(data)
	if e != nil {
		t.Fatal(e)
	}
	actor.Name = "jon"

	data, _ = json.Marshal(actor)
	data = append(data, 1)
	e = service.DeleteActor(data)
	if e == nil {
		t.Fatal(e)
	}

	actor.Name = "jon"

	data, _ = json.Marshal(actor)

	e = service.DeleteActor(data)
	if e != nil {
		t.Fatal(e)
	}

	data, _ = json.Marshal(actor2)

	e = service.DeleteActor(data)
	if e == nil {
		t.Fail()
	}

	// тест 4
	_, e = service.GetListActors()
	if e == nil {
		t.Fatal(e)
	}
	list = 1
	_, e = service.GetListActors()
	if e == nil {
		t.Fatal(e)
	}
	list = 2
	_, e = service.GetListActors()
	if e != nil {
		t.Fatal(e)
	}
	list = 3
	_, e = service.GetListActors()
	if e == nil {
		t.Fatal(e)
	}

	// тест 5

	var film entity.Film
	film.Name = "alive"
	film.Rating = 10
	film.About = "good film"
	film.ReleaseDate = time.Now()

	data, _ = json.Marshal(film)
	e = service.AddFilm(data)
	if e != nil {
		t.Fatal()
	}

	film.Rating = -1

	data, _ = json.Marshal(film)
	e = service.AddFilm(data)
	if e == nil {
		t.Fatal()
	}

	film.Rating = 10
	film.Name = "Warhorse One"
	data, _ = json.Marshal(film)
	e = service.AddFilm(data)
	if e == nil {
		t.Fatal()
	}

	film.Name = "alive"

	data, _ = json.Marshal(film)
	data = append(data, 1)
	e = service.AddFilm(data)
	if e == nil {
		t.Fatal()
	}

	// тест 6
	var film2 entity.Film
	film2.Name = "Warhorse One"
	film2.Rating = 10
	film2.ReleaseDate = time.Now()
	film2.About = "good film"

	var requestfilm entity.RequestEditFilm
	requestfilm.Oldfilm = film2
	requestfilm.NewFilm = film

	data, _ = json.Marshal(requestfilm)
	data = append(data, 1)
	e = service.EditFilm(data)
	if e == nil {
		t.Fatal()
	}

	data, _ = json.Marshal(requestfilm)
	e = service.EditFilm(data)
	if e != nil {
		t.Fatal()
	}
	requestfilm.NewFilm = film2
	data, _ = json.Marshal(requestfilm)
	e = service.EditFilm(data)
	if e == nil {
		t.Fatal()
	}
	requestfilm.Oldfilm = film
	data, _ = json.Marshal(requestfilm)
	e = service.EditFilm(data)
	if e == nil {
		t.Fatal()
	}

	// тест 7

	data, _ = json.Marshal(film2)
	e = service.DeleteFilm(data)
	if e != nil {
		t.Fatal(e)
	}

	data, _ = json.Marshal(film)
	e = service.DeleteFilm(data)
	if e == nil {
		t.Fatal()
	}

	data, _ = json.Marshal(film)
	e = service.DeleteFilm(data)
	if e == nil {
		t.Fatal()
	}

	data, _ = json.Marshal(film)
	data = append(data, 1)
	e = service.DeleteFilm(data)
	if e == nil {
		t.Fatal()
	}

	// тест 8
	_, e = service.GetListFilms(3, -1)
	if e == nil {
		t.Fatal()
	}

	list2 = 1
	_, e = service.GetListFilms(3, -1)
	if e == nil {
		t.Fatal()
	}

	list2 = 2
	_, e = service.GetListFilms(3, -1)
	if e != nil {
		t.Fatal()
	}

	// тест 9
	_, e = service.FindInFilm("error")
	if e == nil {
		t.Fatal()
	}
	_, e = service.FindInFilm("test")
	if e == nil {
		t.Fatal()
	}

	data, e = service.FindInFilm("Warhorse One")
	if e != nil {
		t.Fatal(e)
	}
	var f []entity.Film
	e = json.Unmarshal(data, &f)
	if e != nil {
		t.Fatal(e)
	}
	if f[0].Name != film2.Name {
		t.Fatal()
	}
	var film3 entity.Film
	film3.Id = 1
	film3.Name = "Warhorse One"
	film3.Rating = 10
	film3.ReleaseDate = time.Now()
	film3.About = "good film"
	data, _ = json.Marshal(&film3)
	e = service.DeleteFilm(data)
	if e == nil {
		t.Fatal()
	}
	film3.Id = 0
	film3.Name = "six"
	data, _ = json.Marshal(&film3)
	e = service.DeleteFilm(data)
	if e == nil {
		t.Fatal()
	}
	film3.Rating = -1
	data, _ = json.Marshal(&film3)
	e = service.EditFilm(data)
	if e == nil {
		t.Fatal()
	}

	var actor3 entity.Actor

	actor3.Id = 0
	actor3.Name = "jane"
	actor3.Sex = "12"
	actor3.DataOfBirthday = time.Now()
	request.Oldactor = actor3
	var actor4 entity.Actor

	actor4.Id = 1
	actor4.Name = "jane"
	actor4.Sex = "12"
	actor4.DataOfBirthday = time.Now()
	request.Newactor = actor4

	data, _ = json.Marshal(request)

	e = service.EditActor(data)
	if e == nil {
		t.Fatal()
	}
	data, _ = json.Marshal(actor4)

	e = service.DeleteActor(data)
	if e == nil {
		t.Fatal()
	}

	var film7 entity.Film
	film7.Name = "Warhorse One"
	film7.Id = 1
	film7.Rating = 10
	film7.ReleaseDate = time.Now()
	film7.About = "good film"

	var actor10 entity.Actor

	actor10.Name = "jon"
	actor10.Sex = "male"
	actor10.DataOfBirthday = time.Now()

	var connection entity.RequestEditConnection
	connection.Actor = actor10
	connection.Film = film7

	data, _ = json.Marshal(connection)
	data = append(data, 1)
	e = service.AddConnectionBetweenActorAndFilm(data)
	if e == nil {
		t.Fatal(e)
	}
	data, _ = json.Marshal(connection)
	data = append(data, 1)
	e = service.DeleteConnectionBetweenActorAndFilm(data)
	if e == nil {
		t.Fatal(e)
	}

	data, _ = json.Marshal(connection)

	e = service.AddConnectionBetweenActorAndFilm(data)
	if e != nil {
		t.Fatal(e)
	}

	data, _ = json.Marshal(connection)

	e = service.DeleteConnectionBetweenActorAndFilm(data)
	if e != nil {
		t.Fatal(e)
	}

	connection.Actor.Id = 1

	data, _ = json.Marshal(connection)

	e = service.DeleteConnectionBetweenActorAndFilm(data)
	if e == nil {
		t.Fatal(e)
	}

	data, _ = json.Marshal(connection)

	e = service.AddConnectionBetweenActorAndFilm(data)
	if e == nil {
		t.Fatal(e)
	}

	connection.Actor.Name = "d"

	e = service.DeleteConnectionBetweenActorAndFilm(data)
	if e == nil {
		t.Fatal(e)
	}

	data, _ = json.Marshal(connection)

	e = service.AddConnectionBetweenActorAndFilm(data)
	if e == nil {
		t.Fatal(e)
	}
	connection.Actor.Name = "jon"

	connection.Film.Name = "s"
	data, _ = json.Marshal(connection)

	e = service.DeleteConnectionBetweenActorAndFilm(data)
	if e == nil {
		t.Fatal(e)
	}

	data, _ = json.Marshal(connection)

	e = service.AddConnectionBetweenActorAndFilm(data)
	if e == nil {
		t.Fatal(e)
	}

}

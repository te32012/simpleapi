package auth_test

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"vktestgo2024/internal/auth"
	"vktestgo2024/internal/entity"
)

type Databasemock struct {
}

var userhash string
var adminhash string

func (db *Databasemock) GetUser(ctx context.Context, login string, password string) (*entity.User, error) {
	if "user" == login && userhash == password {
		return &entity.User{Id: 1, Login: "user", Password: "pswrd", Permission: entity.UserPermission}, nil
	}
	if "admin" == login && adminhash == password {
		return &entity.User{Id: 1, Login: "admin", Password: "admin", Permission: entity.AdminPermission}, nil
	}

	return nil, errors.New("тестовая ошибка")
}
func (db *Databasemock) AddActor(ctx context.Context, actor entity.Actor) error  { return nil }
func (db *Databasemock) EditActor(ctx context.Context, actor entity.Actor) error { return nil }
func (db *Databasemock) GetActor(ctx context.Context, actor entity.Actor) (*entity.Actor, error) {
	return nil, nil
}
func (db *Databasemock) DeleteActor(ctx context.Context, actor entity.Actor) error { return nil }
func (db *Databasemock) GetActors(ctx context.Context) ([]entity.Actor, error)     { return nil, nil }
func (db *Databasemock) AddFilm(ctx context.Context, film entity.Film) error       { return nil }
func (db *Databasemock) GetFilm(ctx context.Context, film entity.Film) (*entity.Film, error) {
	return nil, nil
}
func (db *Databasemock) DeleteFilm(ctx context.Context, film entity.Film) error { return nil }
func (db *Databasemock) EditFilm(ctx context.Context, film entity.Film) error   { return nil }
func (db *Databasemock) GetListFilms(ctx context.Context, keySort int, orderSort int) ([]entity.Film, error) {
	return nil, nil
}
func (db *Databasemock) FindInFilm(ctx context.Context, fragment string) ([]entity.Film, error) {
	return nil, nil
}
func (db *Databasemock) GetListFilmByActorId(ctx context.Context, id_actor int) ([]entity.Film, error) {
	return nil, nil
}
func (db *Databasemock) AddActorFilmConnection(ctx context.Context, id_actor, id_film int) error {
	return nil
}
func (db *Databasemock) DeleteActorFilmConnection(ctx context.Context, id_actor, id_film int) error {
	return nil
}
func (c *Databasemock) AddFilmWithActor(ctx context.Context, film entity.Film) error {
	return nil
}

func TestAuthService(t *testing.T) {
	auth := auth.NewAutService(&Databasemock{}, log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile), log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile))
	userhash = fmt.Sprintf("%x", sha256.Sum256([]byte("pswd")))
	adminhash = fmt.Sprintf("%x", sha256.Sum256([]byte("admin")))

	hash, e := auth.LoginUser("user", "pswd")
	if e != nil {
		t.Fatal(e)
		return
	}

	_, e = auth.LoginUser("test", "test")
	if e == nil {
		t.Fatal()
		return
	}
	_, ok := auth.CheckUserIsLoginedAndHasPermission(hash, entity.UserPermission)
	if !ok {
		t.Fatal()
		return
	}
	hash2, e := auth.LoginUser("admin", "admin")
	if e != nil {
		t.Fatal()
		return
	}
	u, ok := auth.CheckUserIsLoginedAndHasPermission(hash2, entity.AdminPermission)
	if !ok || u.Permission == entity.UserPermission {
		t.Fatal()
	}
	_, ok = auth.CheckUserIsLoginedAndHasPermission("123", entity.UserPermission)
	if ok {
		t.Fatal()
		return
	}
	_, ok = auth.CheckUserIsLoginedAndHasPermission("user", -1)
	if ok {
		t.Fatal()
		return
	}

}

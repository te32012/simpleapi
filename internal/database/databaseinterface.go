package database

import (
	"context"
	"vktestgo2024/internal/entity"
)

type DatabaseConnectorInterface interface {
	GetUser(ctx context.Context, login string, password string) (*entity.User, error)
	AddActor(ctx context.Context, actor entity.Actor) error
	EditActor(ctx context.Context, actor entity.Actor) error
	GetActor(ctx context.Context, actor entity.Actor) (*entity.Actor, error)
	DeleteActor(ctx context.Context, actor entity.Actor) error
	GetActors(ctx context.Context) ([]entity.Actor, error)
	AddFilm(ctx context.Context, film entity.Film) error
	GetFilm(ctx context.Context, film entity.Film) (*entity.Film, error)
	DeleteFilm(ctx context.Context, film entity.Film) error
	EditFilm(ctx context.Context, film entity.Film) error
	GetListFilms(ctx context.Context, keySort int, orderSort int) ([]entity.Film, error)
	FindInFilm(ctx context.Context, fragment string) ([]entity.Film, error)
	GetListFilmByActorId(ctx context.Context, id_actor int) ([]entity.Film, error)
	AddActorFilmConnection(ctx context.Context, id_actor, id_film int) error
	DeleteActorFilmConnection(ctx context.Context, id_actor, id_film int) error
	AddFilmWithActor(ctx context.Context, film entity.Film) error
}

package database

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"vktestgo2024/internal/entity"
)

type DatabaseConnector struct {
	Base   *sql.DB
	IsTest bool
}

func NewDatabaseConnector(db *sql.DB, istest bool) *DatabaseConnector {

	return &DatabaseConnector{Base: db, IsTest: istest}
}

func (c *DatabaseConnector) GetUser(ctx context.Context, login string, password string) (*entity.User, error) {
	q := "SELECT id, username, hpassword, permission from users where username=$1 and hpassword=$2;"
	if c.IsTest {
		q = "SELECT|users|id,username,hpassword,permission|username=?,hpassword=?"
	}
	con, err := c.Base.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer con.Close()

	st, err := con.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer st.Close()
	rows, err := st.Query(login, password)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var user entity.User

		err := rows.Scan(&user.Id, &user.Login, &user.Password, &user.Permission)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, errors.New("нужного пользователя нет")
}

func (c *DatabaseConnector) AddActor(ctx context.Context, actor entity.Actor) error {
	q := "INSERT INTO actors (nameActor, sex, dataofbirthday) VALUES (1$, $2, $3);"
	if c.IsTest {
		q = "INSERT|actors|nameActor=?,sex=?,dataofbirthday=?"
	}

	stm, err := c.Base.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stm.Close()

	tx, err := c.Base.Begin()
	if err != nil {
		return err
	}
	res, err := tx.StmtContext(ctx, stm).ExecContext(ctx, actor.Name, actor.Sex, actor.DataOfBirthday)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("не было вставлено ни одной строки в базу")
	}
	return nil
}

// перед этим сделать запрос в базу иначе не проверить существует ли он вообще
// обязательно id актера
func (c *DatabaseConnector) EditActor(ctx context.Context, actor entity.Actor) error {
	q := "UPDATE actors SET nameActor=$1, sex=$2, dataofbirthday=$3 WHERE id=$4;"
	if c.IsTest {
		q = "INSERT|actors|nameActor=?,sex=?,dataofbirthday=?,id=?"
	}

	stm, err := c.Base.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stm.Close()
	tx, err := c.Base.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Stmt(stm).ExecContext(ctx, actor.Name, actor.Sex, actor.DataOfBirthday, actor.Id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("не было изменено ни одной строки в базе")
	}
	return nil
}

func (c *DatabaseConnector) GetActor(ctx context.Context, actor entity.Actor) (*entity.Actor, error) {
	q := "SELECT id, nameActor, sex, dataofbirthday from actors where nameActor=$1 AND dataofbirthday=$2;"
	if c.IsTest {
		q = "SELECT|actors|id,nameActor,sex,dataofbirthday|nameActor=?,dataofbirthday=?"
	}

	stm, err := c.Base.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	r, err := stm.QueryContext(ctx, actor.Name, actor.DataOfBirthday)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	if r.Next() {
		var id int
		var nameActor string
		var sex string
		var dataofbirthday time.Time
		r.Scan(&id, &nameActor, &sex, &dataofbirthday)
		return &entity.Actor{Id: id, Name: nameActor, Sex: sex, DataOfBirthday: dataofbirthday}, nil
	}
	return nil, errors.New("такой пользователь не найден")

}

// обязателен id актера
func (c *DatabaseConnector) DeleteActor(ctx context.Context, actor entity.Actor) error {
	q1 := "DELETE FROM actors WHERE id=$1;"
	if c.IsTest {
		q1 = "INSERT|actors|id=?"
	}
	q2 := "DELETE FROM actors_films WHERE id_actors=$1;"
	if c.IsTest {
		q2 = "INSERT|actors_films|id_actors=?"
	}

	stm, err := c.Base.PrepareContext(ctx, q1)
	if err != nil {
		return err
	}
	defer stm.Close()

	stm2, err := c.Base.PrepareContext(ctx, q2)
	if err != nil {
		return err
	}
	defer stm2.Close()

	tx, err := c.Base.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Stmt(stm).ExecContext(ctx, actor.Id)
	if err != nil {
		return err
	}

	res2, err := tx.Stmt(stm2).ExecContext(ctx, actor.Id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	count2, err2 := res2.RowsAffected()

	if err != nil {
		return err
	}
	if err2 != nil {
		return err2
	}

	if count != 1 || count2 != 1 {
		return errors.New("не было изменено ни одной строки в базе")
	}
	return nil
}

func (c *DatabaseConnector) GetActors(ctx context.Context) ([]entity.Actor, error) {
	q1 := "SELECT id, nameActor, sex, dataofbirthday from actors;"
	if c.IsTest {
		q1 = "SELECT|actors|id,nameActor,sex,dataofbirthday|"
	}

	stm, err := c.Base.PrepareContext(ctx, q1)
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	r, err := stm.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var ans []entity.Actor
	for r.Next() {
		var id int
		var nameActor string
		var sex string
		var dataofbirthday time.Time
		r.Scan(&id, &nameActor, &sex, &dataofbirthday)
		ans = append(ans, entity.Actor{Id: id, Name: nameActor, Sex: sex, DataOfBirthday: dataofbirthday})
	}
	if len(ans) > 0 {
		return ans, nil
	}
	return nil, errors.New("в базе нет пользователей")
}

func (c *DatabaseConnector) AddFilm(ctx context.Context, film entity.Film) error {
	q1 := "INSERT INTO films (nameOfFilm, about, releaseDate, rating) VALUES (1$, $2, $3, $4);"
	if c.IsTest {
		q1 = "INSERT|films|nameOfFilm=?,about=?,releaseDate=?,rating=?"
	}

	stm, err := c.Base.PrepareContext(ctx, q1)
	if err != nil {
		return err
	}
	defer stm.Close()

	tx, err := c.Base.Begin()
	if err != nil {
		return err
	}
	res, err := tx.StmtContext(ctx, stm).ExecContext(ctx, film.Name, film.About, film.ReleaseDate, film.Rating)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("не было вставлено ни одной строки в базу")
	}
	return nil

}

// обязательны только поля имя и дата
func (c *DatabaseConnector) GetFilm(ctx context.Context, film entity.Film) (*entity.Film, error) {
	q := "SELECT id, nameOfFilm, about, releaseDate, rating from films where nameOfFilm=$1 AND releaseDate=$2;"
	if film.Id > 0 {
		q = "SELECT id, nameOfFilm, about, releaseDate, rating from films where id=$1;"
	}
	if c.IsTest {
		if film.Id > 0 {
			q = "SELECT|films|id,nameOfFilm,about,releaseDate,rating|id=?"
		} else {
			q = "SELECT|films|id,nameOfFilm,about,releaseDate,rating|nameOfFilm=?,releaseDate=?"
		}
	}
	stm, err := c.Base.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	var r *sql.Rows
	if film.Id > 0 {
		r, err = stm.QueryContext(ctx, film.Id)
	} else {
		r, err = stm.QueryContext(ctx, film.Name, film.ReleaseDate)
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()
	if r.Next() {
		var ans entity.Film
		r.Scan(&ans.Id, &ans.Name, &ans.About, &ans.ReleaseDate, &ans.Rating)
		return &ans, nil
	}
	return nil, errors.New("такой фильм не найден")

}

// обязателен id фильма
func (c *DatabaseConnector) DeleteFilm(ctx context.Context, film entity.Film) error {
	q1 := "DELETE FROM films WHERE id=$1;"
	q2 := "DELETE FROM actors_films WHERE id_films=$1;"
	q3 := "DELETE FROM fulltextsearch WHERE id_films=$1;"

	if c.IsTest {
		q1 = "INSERT|films|id=?"
		q2 = "INSERT|actors_films|id_films=?"
		q3 = "INSERT|fulltextsearch|id_films=?"
	}

	stm, err := c.Base.PrepareContext(ctx, q1)
	if err != nil {
		return err
	}
	defer stm.Close()
	stm2, err := c.Base.PrepareContext(ctx, q2)
	if err != nil {
		return err
	}
	defer stm2.Close()
	stm3, err := c.Base.PrepareContext(ctx, q3)
	if err != nil {
		return err
	}
	defer stm3.Close()

	tx, err := c.Base.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Stmt(stm).ExecContext(ctx, film.Id)
	tx.Stmt(stm2).ExecContext(ctx, film.Id)
	tx.Stmt(stm3).ExecContext(ctx, film.Id)

	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("не было изменено ни одной строки в базе")
	}
	return nil
}

// обязателен id фильма
func (c *DatabaseConnector) EditFilm(ctx context.Context, film entity.Film) error {
	q1 := "UPDATE films SET nameOfFilm=$1, about=$2, releaseDate=$3, rating=$4 WHERE id=$5"
	if c.IsTest {
		q1 = "INSERT|films|nameOfFilm=?,about=?,releaseDate=?,rating=?,id=?"
	}

	stm, err := c.Base.PrepareContext(ctx, q1)
	if err != nil {
		return err
	}
	defer stm.Close()

	tx, err := c.Base.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Stmt(stm).ExecContext(ctx, film.Name, film.About, film.ReleaseDate, film.Rating, film.Id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("не было изменено ни одной строки в базе")
	}
	return nil
}

// по дефолту запрашивать 3 и -1
func (c *DatabaseConnector) GetListFilms(ctx context.Context, keySort int, orderSort int) ([]entity.Film, error) {
	var q string
	switch keySort {
	case 1:
		if orderSort > 0 {
			q = "SELECT id, nameOfFilm, about, releaseDate, rating from films ORDER BY nameOfFilm;"
		} else {
			q = "SELECT id, nameOfFilm, about, releaseDate, rating from films ORDER BY nameOfFilm DESC;"
		}
	case 2:
		if orderSort > 0 {
			q = "SELECT id, nameOfFilm, about, releaseDate, rating from films ORDER BY releaseDate;"
		} else {
			q = "SELECT id, nameOfFilm, about, releaseDate, rating from films ORDER BY releaseDate DESC;"
		}
	case 3:
		if orderSort > 0 {
			q = "SELECT id, nameOfFilm, about, releaseDate, rating from films ORDER BY rating;"
		} else {
			q = "SELECT id, nameOfFilm, about, releaseDate, rating from films ORDER BY rating DESC;"
		}
	default:
		return nil, errors.New("неправильный ключ")
	}
	if c.IsTest {
		q = "SELECT|films|id,nameOfFilm,about,releaseDate,rating|"
	}

	stm, err := c.Base.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	r, err := stm.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var ans []entity.Film
	for r.Next() {
		var film entity.Film
		r.Scan(&film.Id, &film.Name, &film.About, &film.ReleaseDate, &film.Rating)
		ans = append(ans, film)
	}
	if len(ans) > 0 {
		return ans, nil
	}
	return nil, errors.New("фильмов в базе нет")
}

// ищет только id штуки
func (c *DatabaseConnector) FindInFilm(ctx context.Context, fragment string) ([]entity.Film, error) {
	q := "SELECT id_films FROM fulltextsearch WHERE fulltextsearch.keyworld @@ to_tsquery($1)"
	if c.IsTest {
		q = "SELECT|fulltextsearch|id_films|row=?"
	}
	stm, err := c.Base.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	r, err := stm.QueryContext(ctx, fragment)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var ans []entity.Film
	for r.Next() {
		var id_films int
		r.Scan(&id_films)
		f, err := c.GetFilm(ctx, entity.Film{Id: id_films})
		if err != nil {
			return nil, err
		}
		ans = append(ans, *f)
	}
	if len(ans) > 0 {
		return ans, nil
	}
	return nil, errors.New("совпадений не найдено")
}

func (c *DatabaseConnector) GetListFilmByActorId(ctx context.Context, id_actor int) ([]entity.Film, error) {
	q := "SELECT id_films, id_actors from films where id_actors=$1;"
	if c.IsTest {
		q = "SELECT|actors_films|id_films|id_actors=?"
	}
	stm, err := c.Base.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stm.Close()

	r, err := stm.QueryContext(ctx, id_actor)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var ans []entity.Film
	for r.Next() {
		var id_films int
		r.Scan(&id_films)
		f, err := c.GetFilm(ctx, entity.Film{Id: id_films})
		if err != nil {
			return nil, err
		}
		ans = append(ans, *f)
	}
	if len(ans) > 0 {
		return ans, nil
	}
	return nil, errors.New("такой пользователь не найден")
}

func (c *DatabaseConnector) AddActorFilmConnection(ctx context.Context, id_actor, id_film int) error {
	q := "INSERT INTO films (id_films, id_actors) VALUES (1$, $2);"
	if c.IsTest {
		q = "INSERT|actors_films|id_films=?,id_actors=?"
	}
	stm, err := c.Base.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stm.Close()

	tx, err := c.Base.Begin()
	if err != nil {
		return err
	}
	res, err := tx.StmtContext(ctx, stm).ExecContext(ctx, id_film, id_actor)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("не было вставлено ни одной строки в базу")
	}
	return nil
}

func (c *DatabaseConnector) DeleteActorFilmConnection(ctx context.Context, id_actor, id_film int) error {
	q := "DELETE FROM actors_films WHERE id_films=$1 and id_actors=$2;"
	if c.IsTest {
		q = "INSERT|actors_films|id_films=?,id_actors=?"
	}
	stm, err := c.Base.PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer stm.Close()

	tx, err := c.Base.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Stmt(stm).ExecContext(ctx, id_film, id_actor)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("не было изменено ни одной строки в базе")
	}
	return nil
}

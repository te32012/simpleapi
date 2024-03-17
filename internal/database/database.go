package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"time"
	"vktestgo2024/internal/entity"
)

type DatabaseConnector struct {
	Base   *sql.DB
	IsTest bool
	INFO   *log.Logger
	ERROR  *log.Logger
}

func NewDatabaseConnector(db *sql.DB, istest bool) *DatabaseConnector {
	ans := &DatabaseConnector{Base: db, IsTest: istest}
	ans.INFO = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ans.ERROR = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return ans
}

func (c *DatabaseConnector) GetUser(ctx context.Context, login string, password string) (*entity.User, error) {
	c.INFO.Println("start get user from db")
	q := "SELECT id, username, hpassword, permission from users where username=$1 and hpassword=$2;"
	if c.IsTest {
		q = "SELECT|users|id,username,hpassword,permission|username=?,hpassword=?"
	}
	con, err := c.Base.Conn(ctx)
	if err != nil {
		c.ERROR.Println("error connecting with db")
		return nil, err
	}
	defer con.Close()

	st, err := con.PrepareContext(ctx, q)
	if err != nil {
		c.ERROR.Println("error prepare statment to db")
		return nil, err
	}
	defer st.Close()
	rows, err := st.Query(login, password)
	if err != nil {
		c.ERROR.Println("error query to db")
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var user entity.User
		err := rows.Scan(&user.Id, &user.Login, &user.Password, &user.Permission)
		if err != nil {
			c.ERROR.Println("error scan from db")
			return nil, err
		}
		return &user, nil
	}
	c.ERROR.Println("user with this data not found")
	return nil, errors.New("нужного пользователя нет")
}

func (c *DatabaseConnector) AddActor(ctx context.Context, actor entity.Actor) error {
	q := "INSERT INTO actors (nameActor, sex, dataofbirthday) VALUES ($1, $2, $3);"
	if c.IsTest {
		q = "INSERT|actors|nameActor=?,sex=?,dataofbirthday=?"
	}

	stm, err := c.Base.PrepareContext(ctx, q)
	if err != nil {
		c.ERROR.Println("error prepare statment to db")
		return err
	}
	defer stm.Close()

	tx, err := c.Base.Begin()
	if err != nil {
		c.ERROR.Println("error start transaction in db")
		return err
	}
	res, err := tx.StmtContext(ctx, stm).ExecContext(ctx, actor.Name, actor.Sex, actor.DataOfBirthday)
	if err != nil {
		c.ERROR.Println("error add statment in tx")
		return err
	}
	err = tx.Commit()
	if err != nil {
		c.ERROR.Println("error commit tx")
		tx.Rollback()
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		c.ERROR.Println("error answer from tx")
		return err
	}
	if count != 1 {
		c.ERROR.Println("error answer from tx")
		return errors.New("не было вставлено ни одной строки в базу")
	}
	c.INFO.Println("finish get user from db")
	return nil
}

// перед этим сделать запрос в базу иначе не проверить существует ли он вообще
// обязательно id актера
func (c *DatabaseConnector) EditActor(ctx context.Context, actor entity.Actor) error {
	c.INFO.Println("start edit user in db")

	q := "UPDATE actors SET nameActor=$1, sex=$2, dataofbirthday=$3 WHERE id=$4;"
	if c.IsTest {
		q = "INSERT|actors|nameActor=?,sex=?,dataofbirthday=?,id=?"
	}

	stm, err := c.Base.PrepareContext(ctx, q)
	if err != nil {
		c.ERROR.Println("error prepare stm")
		return err
	}
	defer stm.Close()
	tx, err := c.Base.Begin()
	if err != nil {
		c.ERROR.Println("error start tx")
		return err
	}
	res, err := tx.Stmt(stm).ExecContext(ctx, actor.Name, actor.Sex, actor.DataOfBirthday, actor.Id)
	if err != nil {
		c.ERROR.Println("error add stm to tx")
		return err
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		c.ERROR.Println("error commit tx")
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		c.ERROR.Println("error tx answer")
		return err
	}
	if count != 1 {
		c.ERROR.Println("error tx answer")
		return errors.New("не было изменено ни одной строки в базе")
	}
	c.INFO.Println("finish edit user in db")
	return nil
}

func (c *DatabaseConnector) GetActor(ctx context.Context, actor entity.Actor) (*entity.Actor, error) {
	c.INFO.Println("start get actor from db")
	q := "SELECT id, nameActor, sex, dataofbirthday from actors where nameActor=$1 AND dataofbirthday=$2;"
	if c.IsTest {
		q = "SELECT|actors|id,nameActor,sex,dataofbirthday|nameActor=?,dataofbirthday=?"
	}

	stm, err := c.Base.PrepareContext(ctx, q)
	if err != nil {
		c.ERROR.Println("error prepare stm")
		return nil, err
	}
	defer stm.Close()

	r, err := stm.QueryContext(ctx, actor.Name, actor.DataOfBirthday)
	if err != nil {
		c.ERROR.Println("error query to db")
		return nil, err
	}
	defer r.Close()

	if r.Next() {
		var actor entity.Actor
		r.Scan(&actor.Id, &actor.Name, &actor.Sex, &actor.DataOfBirthday)
		c.INFO.Println(actor)
		c.INFO.Println("get actor from db")
		c.INFO.Println(actor)
		return &actor, nil
	}
	c.ERROR.Println("actor not found in  db")
	return nil, errors.New("такой пользователь не найден")
}

// обязателен id актера
func (c *DatabaseConnector) DeleteActor(ctx context.Context, actor entity.Actor) error {
	c.INFO.Println("start detete actor in db")
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
		c.ERROR.Println("error prepare stm1")
		return err
	}
	defer stm.Close()

	stm2, err := c.Base.PrepareContext(ctx, q2)
	if err != nil {
		c.ERROR.Println("error prepare stm2")
		return err
	}
	defer stm2.Close()

	tx, err := c.Base.Begin()
	if err != nil {
		c.ERROR.Println("error start tx")
		return err
	}
	_, err = tx.Stmt(stm2).ExecContext(ctx, actor.Id)
	res, err := tx.Stmt(stm).ExecContext(ctx, actor.Id)
	if err != nil {
		c.ERROR.Println("error added stm(1,2) in tx")
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		c.ERROR.Println("error commit tx")
		return err
	}

	count, err := res.RowsAffected()

	if err != nil {
		c.ERROR.Println("error answer tx")
		return err
	}
	if count != 1 {
		c.ERROR.Println("error answer tx")
		return errors.New("не было изменено ни одной строки в базе")
	}
	c.INFO.Println("finish detete actor in db")
	return nil
}

func (c *DatabaseConnector) GetActors(ctx context.Context) ([]entity.Actor, error) {
	c.INFO.Println("start get all actors from db")
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
	c.INFO.Println("finish get all actors from db")
	return nil, errors.New("в базе нет пользователей")
}

func (c *DatabaseConnector) AddFilm(ctx context.Context, film entity.Film) error {
	c.INFO.Println("start add film in db")
	q1 := "INSERT INTO films (nameOfFilm, about, releaseDate, rating) VALUES ($1, $2, $3, $4);"
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
	c.INFO.Println("finish add film in db")
	return nil

}

func (c *DatabaseConnector) GetFilm(ctx context.Context, film entity.Film) (*entity.Film, error) {
	c.INFO.Println("finish get film from db")
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
	c.INFO.Println("film in db not found")
	return nil, nil
}

// обязателен id фильма
func (c *DatabaseConnector) DeleteFilm(ctx context.Context, film entity.Film) error {
	c.INFO.Println("delete film from db with dependencies")
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
	tx.Stmt(stm2).ExecContext(ctx, film.Id)
	tx.Stmt(stm3).ExecContext(ctx, film.Id)
	res, err := tx.Stmt(stm).ExecContext(ctx, film.Id)

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
	c.INFO.Println("finish delete film from db with dependencies")
	return nil
}

// обязателен id фильма
func (c *DatabaseConnector) EditFilm(ctx context.Context, film entity.Film) error {
	c.INFO.Println("edit film in db")
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
	c.INFO.Println("finish edit film in db")
	return nil
}

// по дефолту запрашивать 3 и -1
func (c *DatabaseConnector) GetListFilms(ctx context.Context, keySort int, orderSort int) ([]entity.Film, error) {
	c.INFO.Println("start get list films in db")
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
	c.INFO.Println("finish get list films in db")
	return ans, nil
}

// ищет только id штуки
func (c *DatabaseConnector) FindInFilm(ctx context.Context, fragment string) ([]entity.Film, error) {
	c.INFO.Println("start get find by name in db")
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
	c.INFO.Println("finish get find by name in db")
	return nil, errors.New("совпадений не найдено")
}

func (c *DatabaseConnector) GetListFilmByActorId(ctx context.Context, id_actor int) ([]entity.Film, error) {
	c.INFO.Println("start get list film by actor id")
	q := "SELECT id_films from actors_films where id_actors=$1;"
	if c.IsTest {
		q = "SELECT|actors_films|id_films|id_actors=?"
	}
	stm, err := c.Base.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stm.Close()
	c.INFO.Println("осуществление запроса для пользователя ", id_actor)

	r, err := stm.QueryContext(ctx, id_actor)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	c.INFO.Println("поиск фильмов для пользователя", id_actor)

	var ans []entity.Film
	for r.Next() {
		var id_films int
		r.Scan(&id_films)
		f, err := c.GetFilm(ctx, entity.Film{Id: id_films})
		if err != nil {
			return nil, err
		}
		if f != nil {
			ans = append(ans, *f)
		}
	}
	c.INFO.Println("finish get list film by actor id")
	return ans, nil
}

func (c *DatabaseConnector) AddActorFilmConnection(ctx context.Context, id_actor, id_film int) error {
	c.INFO.Println("start add connection between film and actor")
	q := "INSERT INTO actors_films (id_films, id_actors) VALUES ($1, $2);"
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
	c.INFO.Println("finish add connection between film and actor")
	return nil
}

func (c *DatabaseConnector) DeleteActorFilmConnection(ctx context.Context, id_actor, id_film int) error {
	c.INFO.Println("start delete connection between film and actor")
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
	c.INFO.Println("finish delete connection between film and actor")
	return nil
}

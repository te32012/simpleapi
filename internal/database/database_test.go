package database_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"
	"vktestgo2024/internal/database"
	"vktestgo2024/internal/entity"
)

func TestDB(t *testing.T) {
	foo := fakeConnector{name: "foo"}
	sql.Register("fakedb", foo.Driver())
	db, _ := sql.Open("fakedb", "database=root")
	defer db.Close()
	db.Ping()
	conector := database.NewDatabaseConnector(db, true)
	res, err := conector.Base.Exec("CREATE|actors_films|id_films=int32,id_actors=int32")
	if err != nil {
		t.Fatal(err)
	}

	res, err = conector.Base.Exec("CREATE|users|id=int32,username=string,hpassword=string,permission=int32")
	if err != nil {
		t.Fatal(err)
	}
	res, err = conector.Base.Exec("CREATE|fulltextsearch|id_films=int32,row=string")
	if err != nil {
		t.Fatal(err)
	}

	_, err = conector.GetUser(context.Background(), "admin", "admin")
	if err == nil {
		t.Fatal(err)
	}

	res, err = conector.Base.Exec("INSERT|users|id=?,username=user,hpassword=user,permission=?", int32(1), int8(1))
	if err != nil {
		t.Fatal(err)
	}
	count, err := res.RowsAffected()
	if count != 1 {
		t.Fatal(err)
	}
	u, err := conector.GetUser(context.Background(), "user", "user")
	if err != nil || u.Id != 1 || u.Login != "user" || u.Password != "user" || u.Permission != 1 {
		t.Fatal(err)
	}

	res, err = conector.Base.Exec("INSERT|users|id=?,username=admin,hpassword=admin,permission=?", int32(2), int8(2))
	if err != nil {
		t.Fatal(err)
	}
	count, err = res.RowsAffected()
	if count != 1 {
		t.Fatal(err)
	}

	res, err = conector.Base.Exec("CREATE|films|id=int32,nameOfFilm=string,about=string,releaseDate=datetime,rating=int32")
	if err != nil {
		t.Fatal(err)
	}

	res, err = conector.Base.Exec("INSERT|films|id=?,nameOfFilm=alive,about=description,releaseDate=?,rating=?", int32(1), time.Now(), int8(1))
	if err != nil {
		t.Fatal(err)
	}
	count, err = res.RowsAffected()
	if count != 1 {
		t.Fatal(err)
	}

	res, err = conector.Base.Exec("CREATE|actors|id=int32,nameActor=string,sex=string,dataofbirthday=datetime")
	if err != nil {
		t.Fatal(err)
	}
	d0 := time.Now()
	res, err = conector.Base.Exec("INSERT|actors|id=?,nameActor=vasya,sex=male,dataofbirthday=?", int32(1), d0)
	if err != nil {
		t.Fatal(err)
	}
	count, err = res.RowsAffected()
	if count != 1 {
		t.Fatal(err)
	}
	d := time.Now()
	err = conector.AddActor(context.Background(), entity.Actor{Name: "pet", Sex: "male", DataOfBirthday: d})
	if err != nil {
		t.Fatal(err)
	}
	var nameActor string
	var sex string
	var dataofbirthday time.Time
	rows, err := conector.Base.Query("SELECT|actors|nameActor,sex,dataofbirthday|nameActor=?", "pet")
	if err != nil {
		t.Fatal(err)
	}
	if rows.Next() {
		rows.Scan(&nameActor, &sex, &dataofbirthday)
	} else {
		t.Fatal("нет данных")
	}

	if nameActor != "pet" || sex != "male" {
		fmt.Println("пример:", nameActor, " ", sex, dataofbirthday)
		t.Fatal(nameActor, " ", sex)
	}
	err = conector.EditActor(context.Background(), entity.Actor{Id: 10, Name: "pet", Sex: "male", DataOfBirthday: d})
	if err != nil {
		t.Fatal(err)
	}
	var id int
	rows, err = conector.Base.Query("SELECT|actors|id,nameActor,sex,dataofbirthday|id=?", 10)
	if err != nil {
		t.Fatal(err)
	}
	if rows.Next() {
		rows.Scan(&id, &nameActor, &sex, &dataofbirthday)
	} else {
		t.Fatal("нет данных")
	}

	if nameActor != "pet" || sex != "male" || id != 10 || dataofbirthday != d {
		// fmt.Println("пример:", nameActor, " ", sex, dataofbirthday)
		t.Fatal(nameActor, " ", sex, " ", id, " ", dataofbirthday)
	}
	rows, err = conector.Base.Query("SELECT|actors|id,nameActor,sex,dataofbirthday|nameActor=?", "vasya")
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&id, &nameActor, &sex, &dataofbirthday)
		fmt.Println(id, ":", nameActor, ":", sex, ":", dataofbirthday)
	}

	act, err := conector.GetActor(context.Background(), entity.Actor{Name: "vasya", DataOfBirthday: d0})
	if err != nil {
		t.Fatal(err)
	}
	if act.Id != 1 || act.Name != "vasya" || act.Sex != "male" || act.DataOfBirthday != d0 {
		t.Fatal(act.Id, " ", act.Name, " ", act.Sex, " ", act.DataOfBirthday)
	}

	err = conector.DeleteActor(context.Background(), entity.Actor{Id: 11})

	if err != nil {
		t.Fatal(err)
	}
	lst, err := conector.GetActors(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(lst)
	err = conector.AddFilm(context.Background(), entity.Film{Name: "test", About: "about", Rating: 1, ReleaseDate: d0})
	if err != nil {
		t.Fatal(err)
	}
	film, err := conector.GetFilm(context.Background(), entity.Film{Id: 1})
	if err != nil {
		t.Fatal(err)
	}
	if film.Id != 1 && film.Name != "alive" {
		t.Fatal(id, ":", film.Name)
	}
	res, err = conector.Base.Exec("INSERT|films|id=?,nameOfFilm=?,about=description,releaseDate=?,rating=?", int32(5), "test", d, int8(1))
	if err != nil {
		t.Fatal(err)
	}
	count, err = res.RowsAffected()
	if count != 1 {
		t.Fatal(err)
	}
	film, err = conector.GetFilm(context.Background(), entity.Film{Name: "test", ReleaseDate: d})
	if err != nil {
		t.Fatal(err)
	}
	if film.Name != "test" {
		t.Fatal(film)
	}
	film, err = conector.GetFilm(context.Background(), entity.Film{Name: "d", ReleaseDate: d})
	if err == nil {
		t.Fatal(err)
	}
	err = conector.DeleteFilm(context.Background(), entity.Film{Id: 15})
	if err != nil {
		t.Fatal(err)
	}
	err = conector.EditFilm(context.Background(), entity.Film{Id: 30, Name: "fixed", About: "six", Rating: 2, ReleaseDate: d0})
	if err != nil {
		t.Fatal(err)
	}
	films, err := conector.GetListFilms(context.Background(), 3, -1)
	if err != nil {
		t.Fatal(err)
	}
	films, err = conector.GetListFilms(context.Background(), 3, 1)
	if err != nil {
		t.Fatal(err)
	}

	films, err = conector.GetListFilms(context.Background(), 2, -1)
	if err != nil {
		t.Fatal(err)
	}
	films, err = conector.GetListFilms(context.Background(), 2, 1)
	if err != nil {
		t.Fatal(err)
	}

	films, err = conector.GetListFilms(context.Background(), 1, -1)
	if err != nil {
		t.Fatal(err)
	}
	films, err = conector.GetListFilms(context.Background(), 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	films, err = conector.GetListFilms(context.Background(), 0, -1)
	if err == nil {
		t.Fatal(err)
	}
	fmt.Println(films)
	res, err = conector.Base.Exec("INSERT|fulltextsearch|id_films=?,row=test", int32(1))
	if err != nil {
		t.Fatal(err)
	}
	count, err = res.RowsAffected()
	if count != 1 {
		t.Fatal(err)
	}

	films, err = conector.FindInFilm(context.Background(), "test")
	if err != nil {
		t.Fatal(err)
	}
	res, err = conector.Base.Exec("INSERT|films|id=?,nameOfFilm=?,about=description,releaseDate=?,rating=?", int32(42), "test5", d, int8(1))
	if err != nil {
		t.Fatal(err)
	}
	count, err = res.RowsAffected()
	if count != 1 {
		t.Fatal(err)
	}
	res, err = conector.Base.Exec("INSERT|actors|id=?,nameActor=vvv,sex=male,dataofbirthday=?", int32(42), d)
	if err != nil {
		t.Fatal(err)
	}
	count, err = res.RowsAffected()
	if count != 1 {
		t.Fatal(err)
	}
	res, err = conector.Base.Exec("INSERT|actors_films|id_films=?,id_actors=?", int32(42), int32(42))
	if err != nil {
		t.Fatal(err)
	}
	count, err = res.RowsAffected()
	if count != 1 {
		t.Fatal(err)
	}

	films, err = conector.GetListFilmByActorId(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	err = conector.AddActorFilmConnection(context.Background(), 10, 42)
	if err != nil {
		t.Fatal(err)
	}
	err = conector.DeleteActorFilmConnection(context.Background(), 11, 42)
	if err != nil {
		t.Fatal(err)
	}
}

package app

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vktestgo2024/internal/auth"
	"vktestgo2024/internal/database"
	"vktestgo2024/internal/entity"
	"vktestgo2024/internal/midlware"
	"vktestgo2024/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

const uridb = "postgres://root:root@postgres:5432/root"

func Run() {
	tmp1, err := pgxpool.New(context.Background(), uridb)
	db := stdlib.OpenDBFromPool(tmp1)
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(time.Second * 10)
	connector := database.NewDatabaseConnector(db, false)
	Info := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	auth := auth.NewAutService(connector, Info, Error)
	service := service.NewService(auth, connector)

	router := midlware.NewRouter("vktestgo", "2024", service)
	tmp := fmt.Sprintf("%x", sha256.Sum256(([]byte("user"))))
	fmt.Println(tmp)
	db.Exec("INSERT INTO users (username, hpassword, permission) VALUES ($1, $2, $3);", "user", tmp, entity.UserPermission)
	tmp = fmt.Sprintf("%x", sha256.Sum256([]byte("admin")))
	fmt.Println(tmp)
	db.Exec("INSERT INTO users (username, hpassword, permission) VALUES ($1, $2, $3);", "admin", tmp, entity.AdminPermission)
	db2, _ := sql.Open("pgx", uridb)
	connector2 := database.NewDatabaseConnector(db2, false)
	go func() {
		for {
			var next chan int = make(chan int)
			go connector2.Updatekeywords(next)
			<-next
			fmt.Println("STOPPED UPDATE")
			time.Sleep(time.Second * 30)
		}
	}()
	go router.Lisen()

	exit := make(chan os.Signal, 2)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

}

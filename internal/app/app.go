package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"vktestgo2024/internal/auth"
	"vktestgo2024/internal/database"
	"vktestgo2024/internal/entity"
	"vktestgo2024/internal/midlware"
	"vktestgo2024/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const uridb = "postgres://root:pass@postgres:5432/root"

func Run() {
	db, err := sql.Open("pgx", uridb)
	if err != nil {
		log.Fatalln(err)
	}
	connector := database.NewDatabaseConnector(db, false)
	Info := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	auth := auth.NewAutService(connector, Info, Error)
	service := service.NewService(auth, connector)

	router := midlware.NewRouter("vktestgo", "2024", service)
	tmp := fmt.Sprintf("%x", auth.Hasher.Sum([]byte("user")))
	fmt.Println(tmp)
	db.Exec("INSERT INTO users (username, hpassword, permission) VALUES ($1, $2, $3);", "user", tmp, entity.UserPermission)
	tmp = fmt.Sprintf("%x", auth.Hasher.Sum([]byte("admin")))
	fmt.Println(tmp)
	db.Exec("INSERT INTO users (username, hpassword, permission) VALUES ($1, $2, $3);", "admin", tmp, entity.AdminPermission)

	go router.Lisen()

	exit := make(chan os.Signal, 2)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	<-exit

}

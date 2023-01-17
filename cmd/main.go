package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var (
	PortHandler = ""
	Handler     = ""
	PathBD      = ""
	db          *sql.DB
)

func Init() {
	var filearray [6]string
	file, err := os.Open("config.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {

		filearray[i] = scanner.Text()
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	PortHandler = filearray[0]
	Handler = filearray[1]
	PathBD = filearray[2]
}

func main() {
	Init()
	fmt.Println("Запущенно")
	r := mux.NewRouter()
	r.HandleFunc("/marlo/admin/adduser", AuthorizationAdmin(AddUser)).Methods("Get")                             // добавление тестовых пользователей
	r.HandleFunc("/marlo/admin/getdocpattern", AuthorizationAdmin(Getdockspattern)).Methods("Get")               // получения списка шаблонов
	r.HandleFunc("/marlo/admin/adddocpattern", AuthorizationAdmin(Adddockpattern)).Methods("Post")               // создание шаблона
	r.HandleFunc("/marlo/admin/deletedocpattern/{id}/", AuthorizationAdmin(Deletedockpattern)).Methods("Delete") // удаление шаблона
	r.HandleFunc("/marlo/admin/searchdockspattern", AuthorizationAdmin(Searchdockspattern)).Methods("Get")       // поиск шаблонов
	r.HandleFunc("/marlo/admin/updatedockpattern/{id}/", AuthorizationAdmin(Updatedockpattern)).Methods("Put")   //обновление шаблона
	r.HandleFunc("/marlo/admin/getdockstext", AuthorizationAdmin(Getdockstext)).Methods("Get")
	r.HandleFunc("/marlo/admin/getdockstextbydocid/{id}/", AuthorizationAdmin(Getdockstextbydocid)).Methods("Get")
	r.HandleFunc("/marlo/admin/getdockstextbyid/{id}/", AuthorizationAdmin(Getdockstextbyid)).Methods("Get")
	r.HandleFunc("/marlo/admin/getdockstextactyality/{id}/", AuthorizationAdmin(Getdockstextactyality)).Methods("Get")
	r.HandleFunc("/marlo/admin/adddockstextactyality/{id}/", AuthorizationAdmin(Adddockstextactyality)).Methods("Post")
	log.Fatal(http.ListenAndServe(PortHandler, r))

}

func Sqlconnectionmarlo(namebd string) {
	//"root:1234@tcp(localhost:3306)/admin"
	cfg := mysql.Config{
		User:   "root",
		Passwd: "1234",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: namebd,
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

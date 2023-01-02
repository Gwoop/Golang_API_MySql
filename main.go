package main

import (
	"AdminSimpleApi/Structs"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ddosify/go-faker/faker"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
)

func main() {

	fmt.Println("Запущенно")
	r := mux.NewRouter()
	//r.HandleFunc("/marlo/admin/authorization", AuthorizationAdmin).Methods("Get")
	r.HandleFunc("/marlo/admin/adduser", AuthorizationAdmin(AddUser)).Methods("Get")

	log.Fatal(http.ListenAndServe(":8000", r))

}

//Функция хэширования пароля
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

//хэндлел для атворизации
func AuthorizationAdmin(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, password, ok := r.BasicAuth()
		if !ok {
			(w).WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Ошибка обработки")
			return
		}
		db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/admin")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		rows, err := db.Query("select * from admin.aunt")

		for rows.Next() {
			p := admin{}
			erro := rows.Scan(&p.id, &p.login, &p.password)
			if erro != nil {
				fmt.Println(erro)
				continue
			}
			if login != p.login || password != p.password {
				(w).WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode("Ошибка данных")
				return
			}
		}
		(w).WriteHeader(http.StatusOK)
		handler.ServeHTTP(w, r)
	}
}

type admin struct {
	id       int
	login    string
	password string
}

//хэндлер для добавления тестовых пользователей
func AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userst := Structs.User{}
	respusers := Structs.ResponsesUser{}

	var tokenstruct Structs.Token
	_ = json.NewDecoder(r.Body).Decode(&tokenstruct)

	faker := faker.NewFaker()                       //создание экземпляра объекта faker для подставных данных
	userst.Name = faker.RandomPersonFirstName()     //генерация подставного имяни
	userst.Lastname = faker.RandomPersonFirstName() //генерация подставной фамилии
	userst.Sex = rand.Intn(11-0) + 0                //генерация подставного пола (значение от 0 до 11 включительно)
	userst.Birdh = "2008/10/23"                     //Подставление данной даты рождения
	userst.Tel = faker.RandomPhoneNumber()          //генерация подставного номера телефона (логина)
	userst.Chatid = rand.Intn(11-0) + 0             //генерация подставного id чата (значение от 0 до 11 включительно)
	userst.Email = faker.RandomEmail()              //генерация подставной почты
	userst.Password = faker.RandomPassword()        //генерация подставного пороля (здесь он в чистом ввиде без хэширования)

	db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/marlo") //строка подключения к бд (root - пользователь, 1234- пароль от бд,
	// tcp(localhost:3306) - тип подлючения с путём поключения), marlo - название бд
	//обработка ошибки  при подключении к бд
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//запрос для бд на добавления тестового пользователя (пароль хешируется под MD5 с помощью функции GetMD5Hash)
	result, err := db.Exec("insert into marlo.users (name, lastname, sex,birdh,tel,chatid,email,password) values ( ?, ?, ?, ?, ?, ?, ?, ?)",
		userst.Name, userst.Lastname, userst.Sex, userst.Birdh, userst.Tel, userst.Chatid, userst.Email, GetMD5Hash(userst.Password))
	if err != nil {
		panic(err)
	}

	respusers.Id, _ = result.LastInsertId() // нужно обработать ошибку
	respusers.Login = userst.Tel
	respusers.Password = userst.Password

	json.NewEncoder(w).Encode(&respusers)

	fmt.Println(result.LastInsertId()) // id добавленного объекта
	fmt.Println(result.RowsAffected()) // количество затронутых строк
}

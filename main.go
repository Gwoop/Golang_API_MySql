package main

import (
	"AdminSimpleApi/Structs"
	"bufio"
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
	"os"
)

var (
	PortHandler = ""
	Handler     = ""
	PathBD      = ""
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
	r.HandleFunc(Handler, AuthorizationAdmin(AddUser)).Methods("Get")
	log.Fatal(http.ListenAndServe(PortHandler, r))
}

//Функция хэширования пароля
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

//хэндлел для авторизации
func AuthorizationAdmin(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, password, ok := r.BasicAuth() //инициализация базовой авторизации
		if !ok {
			(w).WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Ошибка обработки сессии")
			return
		}
		db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/admin") //строка подлючения к бд где хранятся логин и пароль админа
		if err != nil {
			panic(err)
		}
		defer db.Close()

		//проверка логина и пароля администратора из базы данных
		rows, err := db.Query("select * from admin.aunt")
		for rows.Next() {
			p := Structs.Admin{}
			erro := rows.Scan(&p.Id, &p.Login, &p.Password)
			if erro != nil {
				fmt.Println(erro)
				continue
			}
			if login != p.Login || password != p.Password {
				(w).WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode("Ошибка ввода данных (логин или пароль не верны)")
				return
			}
		}
		(w).WriteHeader(http.StatusOK)
		handler.ServeHTTP(w, r)
	}
}

//хэндлер для добавления тестовых пользователей
func AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//userst := Structs.User{}
	respusers := Structs.ResponsesUser{} //структура для ответа пользователю

	var tokenstruct Structs.Token
	_ = json.NewDecoder(r.Body).Decode(&tokenstruct) //получение токена (не используется в данныйй момент в программе)

	userst := FakeData() //подставление случаных данных для создания тестового пользователя

	db, err := sql.Open("mysql", PathBD) //строка подключения к бд (root - пользователь, 1234- пароль от бд,
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

	respusers.Id, _ = result.LastInsertId() // id созданного пользователя
	respusers.Login = userst.Tel            // телефон (логин) созданного пользователя
	respusers.Password = userst.Password    // пароль (без шифрования) созданного пользователя

	json.NewEncoder(w).Encode(&respusers) // отправка ответа пользователю

	fmt.Println(result.LastInsertId()) // id добавленного объекта
	fmt.Println(result.RowsAffected()) // количество затронутых строк
}

func FakeData() Structs.User {
	userst := Structs.User{}
	faker := faker.NewFaker()                       //создание экземпляра объекта faker для подставных данных
	userst.Name = faker.RandomPersonFirstName()     //генерация подставного имяни
	userst.Lastname = faker.RandomPersonFirstName() //генерация подставной фамилии
	userst.Sex = rand.Intn(11-0) + 0                //генерация подставного пола (значение от 0 до 11 включительно)
	userst.Birdh = "2008/10/23"                     //Подставление данной даты рождения
	userst.Tel = faker.RandomPhoneNumber()          //генерация подставного номера телефона (логина)
	userst.Chatid = rand.Intn(11-0) + 0             //генерация подставного id чата (значение от 0 до 11 включительно)
	userst.Email = faker.RandomEmail()              //генерация подставной почты
	userst.Password = faker.RandomPassword()        //генерация подставного пороля (здесь он в чистом ввиде без хэширования)
	return userst
}

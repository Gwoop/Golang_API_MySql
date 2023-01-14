package main

import (
	"AdminSimpleApi/Structs"
	"AdminSimpleApi/cmd/security"
	"encoding/json"
	"fmt"
	"net/http"
)

//хэндлел для авторизации
func AuthorizationAdmin(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, password, ok := r.BasicAuth() //инициализация базовой авторизации
		if !ok {
			(w).WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Ошибка обработки сессии")
			return
		}
		//db, err := sql.Open("mysql", "root:1234@tcp(localhost:3306)/admin") //строка подлючения к бд где хранятся логин и пароль админа
		Sqlconnectionmarlo("admin")
		var err error
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
	var tokenstruct Structs.Token
	_ = json.NewDecoder(r.Body).Decode(&tokenstruct) //получение токена (не используется в данныйй момент в программе)

	userst := FakeData() //подставление случаных данных для создания тестового пользователя

	Sqlconnectionmarlo("marlo")
	var err error
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//запрос для бд на добавления тестового пользователя (пароль хешируется под MD5 с помощью функции GetMD5Hash)
	result, err := db.Exec("insert into marlo.users (name, lastname, sex,birdh,tel,chatid,email,password) values ( ?, ?, ?, ?, ?, ?, ?, ?)",
		userst.Name, userst.Lastname, userst.Sex, userst.Birdh, userst.Tel, userst.Chatid, userst.Email, security.GetMD5Hash(userst.Password))
	if err != nil {
		panic(err)
	}
	var id, _ = result.LastInsertId()

	json.NewEncoder(w).Encode(Structs.ResponsesUser{Id: id, Login: userst.Tel, Password: userst.Password}) // отправка ответа пользователю
	fmt.Println(result.LastInsertId())                                                                     // id добавленного объекта
	fmt.Println(result.RowsAffected())                                                                     // количество затронутых строк
}

func Getdockspattern(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	Sqlconnectionmarlo("marlo")
	rows, err := db.Query("SELECT * FROM document")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		doc := Structs.ResponsesDockpattern{}
		rows.Scan(&doc.Id, &doc.Name, &doc.Description, &doc.Uuid, &doc.Create_date)
		json.NewEncoder(w).Encode(&doc)
	}
}

func Adddockpattern(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Structs.ResponsesSytem{}
	var requestdockpattern Structs.RequestDockpattern
	_ = json.NewDecoder(r.Body).Decode(&requestdockpattern)

	fmt.Println(requestdockpattern.Name)
	fmt.Println(requestdockpattern.Description)
	fmt.Println(requestdockpattern.Uuid)
	Sqlconnectionmarlo("marlo")
	var err error
	defer db.Close()

	_, err = db.Query("insert into marlo.document (name, description, uuid) values (?,?,?)", requestdockpattern.Name, requestdockpattern.Description, requestdockpattern.Uuid)
	fmt.Println(err)
	if err != nil {
		response.Responses = "Ошибка данных " + err.Error()
		json.NewEncoder(w).Encode(&response)
		return
	}
	defer db.Close()
	response.Responses = "Данные успешно добавлены"
	json.NewEncoder(w).Encode(&response)
}

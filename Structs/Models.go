package Structs

//СТруктура для создания тестового пользователя
type User struct {
	Name     string //Имя
	Lastname string //Фамилия
	Sex      int    //Пол
	Birdh    string //Дата рождения
	Tel      string //номер телефона (логин)
	Chatid   int    //Внешний ключ от таблицы чата
	Email    string //Электронная почта
	Password string //Пароль
}

type Token struct {
	Token string `json:"token"`
}

//Структура ответа
type ResponsesUser struct {
	Id       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

//Структура для авторизации админа
type Admin struct {
	Id       int
	Login    string
	Password string
}

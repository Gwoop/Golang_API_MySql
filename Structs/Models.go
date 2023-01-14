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

//Структура ответа для создания пользоватлей
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

//Структура ответа для просмотра всех шаблонов документов
type ResponsesDockpattern struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Uuid        string `json:"uuid"`
	Create_date string `json:"create_date"`
}

//Структура запроса для создания шаблона документа
type RequestDockpattern struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Uuid        string `json:"uuid"`
}

type RequestDockid struct {
	Id string `json:"id"`
}

type RequestDockuuid struct {
	Uuid string `json:"uuid"`
}

type ResponsesSytem struct {
	Responses string `json:"responses"`
}

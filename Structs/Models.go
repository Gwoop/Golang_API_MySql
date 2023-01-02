package Structs

type User struct {
	Name     string //`json:"name"`
	Lastname string //`json:"lastname"`
	Sex      int    //`json:"sex"`
	Birdh    string //`json:"birdh"`
	Tel      string //`json:"tel"`
	Chatid   int    //`json:"chatid"`
	Email    string //`json:"email"`
	Password string //`json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type ResponsesUser struct {
	Id       int64  `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

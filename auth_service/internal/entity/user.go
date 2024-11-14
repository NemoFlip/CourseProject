package entity

type User struct {
	ID       string `json:"id" example:"1234"`
	Username string `json:"username" example:"user_1234"`
	Email    string `json:"email" example:"user_1234@gmail.com"`
	Phone    string `json:"phone" example:"89178298123"`
	Password string `json:"password" example:"pass_1234"`
}

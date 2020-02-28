package vm

type User struct {
	Username    string `json:"user_name"`
	Password    string `json:"user_password"`
	Email       string `json:"user_email"`
	Realname    string `json:"user_realname"`
	Comment     string `json:"user_comment"`
	SystemAdmin int    `json:"user_system_admin"`
}

package vm

type ChangePassword struct {
	OldPassword string `json:"user_password_old"`
	NewPassword string `json:"user_password_new"`
}

package users_services

type UsrLoginparam struct {
	ID       int    `json:"id"`
	UserId   string `json:"userid"`
	Password string `json:"password"`
}
type AuthinticationParam struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}
type User struct {
	ID       int    `json:"id"`
	Username string `json:"userId"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// type UsrClaim struct {
// 	UserID   int      `json:"user_id"`
// 	AuthList []string `json:"auth_list"`
// }

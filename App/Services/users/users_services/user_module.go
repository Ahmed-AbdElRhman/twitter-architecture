package users_services

type UsrLoginparam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	AuthList []string `json:"auth_list"`
}

type UsrRegisterparam struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// type UsrClaim struct {
// 	UserID   int      `json:"user_id"`
// 	AuthList []string `json:"auth_list"`
// }

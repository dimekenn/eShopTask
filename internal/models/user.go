package models

type User struct {
	Id        int      `json:"id"`
	Username  string   `json:"username"`
	Password  string   `json:"password,omitempty"`
	Roles     []int    `json:"roles,omitempty"`
	RolesName []string `json:"roles_name"`
}

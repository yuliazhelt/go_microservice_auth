package models

type User struct {
    Login string `mapstructure:"login"`
	Password string `mapstructure:"password"`
	Role string `mapstructure:"role"`
}
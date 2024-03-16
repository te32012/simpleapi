package auth

import "vktestgo2024/internal/entity"

type AuthServiceInterface interface {
	LoginUser(login string, password string) (string, error)
	CheckUserIsLoginedAndHasPermission(key string, operation int) (*entity.User, bool)
}

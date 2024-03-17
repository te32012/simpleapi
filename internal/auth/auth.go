package auth

import (
	"context"
	"crypto"
	"errors"
	"fmt"
	"hash"
	"log"
	"time"
	"vktestgo2024/internal/database"
	"vktestgo2024/internal/entity"
)

type AuthService struct {
	Session           map[string]*entity.User
	DatabaseConnector database.DatabaseConnectorInterface
	INFO              *log.Logger
	ERROR             *log.Logger
	Hasher            hash.Hash
}

func NewAutService(base database.DatabaseConnectorInterface, info *log.Logger, err *log.Logger) *AuthService {
	return &AuthService{Session: make(map[string]*entity.User), DatabaseConnector: base, Hasher: crypto.SHA256.New(), INFO: info, ERROR: err}
}

func (s *AuthService) LoginUser(login string, password string) (string, error) {
	tmp := fmt.Sprintf("%x", s.Hasher.Sum([]byte(password)))
	// s.INFO.Println(tmp)
	user, err := s.DatabaseConnector.GetUser(context.Background(), login, tmp)
	if err != nil {
		s.ERROR.Println("ошибка получения юзера")
		return "", err
	}
	var hashstr string = ""
	var ok = true
	for ok {
		var hash = crypto.SHA256.New().Sum([]byte(user.Login + time.Now().GoString()))
		hashstr = fmt.Sprintf("%x", hash[:16])
		_, ok = s.Session[hashstr]
	}

	s.Session[hashstr] = user
	return hashstr, nil
}

func (s *AuthService) CheckUserIsLoginedAndHasPermission(key string, operation int) (*entity.User, bool) {
	u, e := s.userIdentificationBySession(key)
	s.ERROR.Println("ошибка идентификации юзера")
	if e != nil {
		return nil, false
	}
	return u, s.checkingPermissionToPerformAnOperation(u, operation)
}

func (s *AuthService) userIdentificationBySession(key string) (*entity.User, error) {
	u, ok := s.Session[key]
	if !ok {
		return nil, errors.New("пользователя для такой сессии не существует")
	}
	return u, nil
}
func (s *AuthService) checkingPermissionToPerformAnOperation(user *entity.User, operation int) bool {
	switch user.Permission {
	case 2:
		return true
	case 1:
		return user.Permission == operation
	default:
		return false
	}
}

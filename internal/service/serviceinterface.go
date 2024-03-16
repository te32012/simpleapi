package service

type ServiceInterface interface {
	Login(login, password string) (string, error)
	CheckUserIsLoginedAndHasPermission(key string, operation int) bool
	AddActor(data []byte) error
	EditActor(data []byte) error
	DeleteActor(data []byte) error
	GetListActors() ([]byte, error)
	// добавляем единственного актера без id и без films
	AddFilm(data []byte) error
	// редактируем единственного актера без id и без films
	EditFilm(data []byte) error
	GetListFilms(keySort int, orderSort int) ([]byte, error)
	DeleteFilm(data []byte) error
	FindInFilm(segment string) ([]byte, error)
	AddConnectionBetweenActorAndFilm(data []byte) error
	DeleteConnectionBetweenActorAndFilm(data []byte) error
}

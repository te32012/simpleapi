package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"vktestgo2024/internal/auth"
	"vktestgo2024/internal/database"
	"vktestgo2024/internal/entity"
)

type Service struct {
	AuthService       auth.AuthServiceInterface
	DatabaseConnector database.DatabaseConnectorInterface
	Info              *log.Logger
	Error             *log.Logger
}

func NewService(auth auth.AuthServiceInterface, base database.DatabaseConnectorInterface) *Service {
	ans := &Service{}
	ans.AuthService = auth
	ans.DatabaseConnector = base
	ans.Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ans.Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return ans
}

func (s *Service) Login(login, password string) (string, error) {
	s.Info.Println(login, password)
	return s.AuthService.LoginUser(login, password)
}

func (s *Service) CheckUserIsLoginedAndHasPermission(key string, operation int) bool {
	err, ans := s.AuthService.CheckUserIsLoginedAndHasPermission(key, operation)
	s.Error.Println(err)
	return ans
}

func (s *Service) AddActor(data []byte) error {
	s.Info.Println("start adding one actor without films and id")
	var actor entity.Actor
	err := json.Unmarshal(data, &actor)
	if err != nil {
		s.Error.Println("error unmarshaling actor")
		return err
	}
	if e := s.addoneactor(actor); e != nil {
		s.Error.Println("error adding actor without films and id")
		return e
	}
	s.Info.Println("actor without films and id added")
	return nil
}

func (s *Service) EditActor(data []byte) error {
	s.Info.Println("start edit one actor without films and id")
	var request entity.RequestEditActor
	err := json.Unmarshal(data, &request)
	if err != nil {
		s.Error.Println("error unmarshaling request edit actor without films and id")
		return err
	}
	if e := s.editoneactor(request); e != nil {
		s.Error.Println("error edit actor without films and id")
		return e
	}
	s.Info.Println("finished edit one actor without films and id")
	return nil
}

func (s *Service) DeleteActor(data []byte) error {
	s.Info.Println("start delete one actor without films and id")
	var actor entity.Actor
	err := json.Unmarshal(data, &actor)
	if err != nil {
		s.Error.Println("error unmarshal actors")
		return err
	}
	if e := s.deleteoneactor(actor); e != nil {
		s.Error.Println("error delete one actor without films and id")
		return e
	}
	s.Info.Println("finised deleting one actor without films and id")
	return nil
}

func (s *Service) GetListActors() ([]byte, error) {
	s.Info.Println("start get list actors")
	actors, err := s.DatabaseConnector.GetActors(context.Background())
	if err != nil {
		return nil, err
	}
	s.Info.Println("geted actors without films")
	s.Info.Println(actors)
	s.Info.Println("start getting films for one actor")

	for i := 0; i < len(actors); i++ {
		s.Info.Println("get film by actor ", actors[i].Id)
		films, err := s.DatabaseConnector.GetListFilmByActorId(context.Background(), actors[i].Id)
		if err != nil {
			s.Error.Println("error get film by actor ", actors[i].Id)
			return nil, err
		}
		actors[i].Films = films
		s.Info.Println("films for actor id ", actors[i].Id, " addied")

	}
	s.Info.Println("finished getting films for one actor")
	if len(actors) == 0 {
		s.Error.Println("actors not found")
		return nil, errors.New("нет актеров")
	}
	ans, err := json.Marshal(actors)
	if err != nil {
		s.Error.Println("actors marshaling error")
		return nil, err
	}
	s.Info.Println("finished get list actors")
	return ans, nil
}

func (s *Service) AddFilm(data []byte) error {
	s.Info.Println("start adding one film without actors and id")
	var film entity.Film
	err := json.Unmarshal(data, &film)
	if err != nil {
		s.Error.Println("error unmarshal film")
		return err
	}
	if e := s.addfilm(film); e != nil {
		s.Error.Println("error adding one film without actors and id")
		return e
	}
	s.Info.Println("finised adding one film without actors and id")
	return nil
}

func (s *Service) EditFilm(data []byte) error {
	s.Info.Println("start edit one film without actors and id")
	var request entity.RequestEditFilm
	err := json.Unmarshal(data, &request)
	if err != nil {
		s.Error.Println("error unmarshal films")
		return err
	}
	if e := s.editfilm(request); e != nil {
		s.Info.Println("error edit one film without actors and id")
		return e
	}
	s.Info.Println("finish edit one film without actors and id")
	return nil
}

func (s *Service) DeleteFilm(data []byte) error {
	s.Info.Println("start delete one film without actors and id")
	var film entity.Film
	err := json.Unmarshal(data, &film)
	if err != nil {
		s.Error.Println("error unmarshal film")
		return err
	}
	if e := s.deltefilm(film); e != nil {
		s.Info.Println("error delete one film without actors and id")
		return e
	}
	s.Info.Println("finished delete one film without actors and id")
	return nil
}

func (s *Service) GetListFilms(keySort int, orderSort int) ([]byte, error) {
	s.Info.Println("start getlistfilms")
	films, err := s.DatabaseConnector.GetListFilms(context.Background(), keySort, orderSort)
	if err != nil {
		s.Error.Println("error getlisfilms")
		return nil, err
	}
	if len(films) == 0 {
		s.Error.Println("films not found")
		return nil, errors.New("нет фильмов")
	}
	s.Info.Println("geted data")
	var ans []byte
	ans, err = json.Marshal(films)
	if err != nil {
		s.Error.Println("error marshal films")
		return nil, err
	}
	s.Info.Println("finished getlistfilms")
	return ans, nil
}

func (s *Service) FindInFilm(segment string) ([]byte, error) {
	s.Info.Println("start findinfilm")
	s.Info.Println(segment)
	films, err := s.DatabaseConnector.FindInFilm(context.Background(), segment)
	if err != nil {
		s.Error.Println("error find in film")
		return nil, err
	}

	if len(films) == 0 {
		s.Error.Println("films not found")
		return nil, errors.New("фильмов не найдено")
	}
	s.Info.Println("fims getting")

	data, err := json.Marshal(films)
	// fmt.Printf("%s", data)
	if err != nil {
		s.Error.Println("error marshaling")
		return nil, err
	}
	s.Info.Println("end findinfilm")
	return data, nil
}

func (s *Service) addoneactor(actor entity.Actor) error {
	s.Info.Println("start adding one actor")
	s.Info.Println(actor)
	if !(actor.Sex == entity.SexMale || actor.Sex == entity.SexFemale) || actor.Id != 0 || actor.Films != nil || len(actor.Films) != 0 {
		s.Error.Println("error validatate data actor")
		return errors.New("error validatate data actor")
	}

	err := s.DatabaseConnector.AddActor(context.Background(), actor)
	if err != nil {
		s.Error.Println("error adding one actor in database")
		return err
	}
	s.Info.Println("finished adding one actor")
	return nil
}

func (s *Service) editoneactor(request entity.RequestEditActor) error {
	s.Info.Println("start edit one actor")
	s.Info.Println(request)
	if request.Oldactor.Id != 0 || request.Newactor.Id != 0 || request.Oldactor.Films != nil ||
		request.Newactor.Films != nil || len(request.Oldactor.Films) != 0 || len(request.Newactor.Films) != 0 ||
		!(request.Newactor.Sex == entity.SexMale || request.Newactor.Sex == entity.SexFemale || request.Newactor.Sex == "") {
		s.Error.Println("ivalide data request for edit one actor")
		return errors.New("не валидные данные для редактирования")
	}

	s.Info.Println("getting data about actor from base")
	u, e := s.DatabaseConnector.GetActor(context.Background(), request.Oldactor)
	if e != nil {
		s.Error.Println("error geting data about edit old actor")
		return e
	}
	s.Info.Println("finised geting data about edit old actor")
	s.Info.Println(*u)
	request.Oldactor = *u
	if request.Newactor.Name != "" {
		s.Info.Println("editing name old actor")
		request.Oldactor.Name = request.Newactor.Name
	}
	if request.Newactor.Sex != "" {
		s.Info.Println("editing sex old actor")
		request.Oldactor.Sex = request.Newactor.Sex
	}
	if !request.Newactor.DataOfBirthday.IsZero() {
		s.Info.Println("editing dataofbirthday old actor")
		request.Oldactor.DataOfBirthday = request.Newactor.DataOfBirthday
	}
	s.Info.Println("start edit data about edit old actor")
	e = s.DatabaseConnector.EditActor(context.Background(), request.Oldactor)
	if e != nil {
		s.Error.Println("error edit data about edit old actor")
		return e
	}
	s.Info.Println("finish edit one actor")
	return nil
}

func (s *Service) deleteoneactor(actor entity.Actor) error {
	s.Info.Println("start gettting one actor without films and id when we need to delete")
	if actor.Id != 0 {
		s.Error.Println("ivalide data request for edit one actor")
		return errors.New("не валидные данные для удаления")
	}
	u, e := s.DatabaseConnector.GetActor(context.Background(), actor)
	if e != nil {
		s.Error.Println("error get data about actor")
		return e
	}
	s.Info.Println(*u)
	s.Info.Println("start delete one actor")
	e = s.DatabaseConnector.DeleteActor(context.Background(), *u)
	if e != nil {
		s.Error.Println("error deleting actor")
		return e
	}
	return nil
}

func (s *Service) addfilm(film entity.Film) error {
	s.Info.Println("start adding film")
	s.Info.Println(film)
	var valid bool = true
	valid = valid && len(film.Name) > 1 && len(film.Name) < 150
	valid = valid && len(film.About) < 1000
	valid = valid && film.Rating >= 0 && film.Rating <= 10
	valid = valid && film.Id == 0
	valid = valid && (film.Actors == nil || len(film.Actors) == 0)
	if !valid {
		s.Error.Println("invalide data about film")
		return errors.New("полученные данные про фильм не корректные")
	}
	s.Info.Println("start adding film")
	e := s.DatabaseConnector.AddFilm(context.Background(), film)
	if e != nil {
		s.Error.Println("error adding film")
		return e
	}
	s.Info.Println("finised adding film")
	return nil
}

func (s *Service) editfilm(request entity.RequestEditFilm) error {
	s.Info.Println("start editing one filmrequest")
	s.Info.Println(request)

	if len(request.NewFilm.Name) >= 150 || len(request.NewFilm.About) >= 1000 || request.NewFilm.ReleaseDate.IsZero() || (request.NewFilm.Rating < 0 && request.NewFilm.Rating > 10) {
		s.Error.Println("invalide data request")
		return errors.New("некорректные данные")
	}

	s.Info.Println("get films from base")
	f, e := s.DatabaseConnector.GetFilm(context.Background(), request.Oldfilm)
	if e != nil {
		s.Error.Println("error getting one film")
		return e
	}
	s.Info.Println("finish getting film without actors and id")
	s.Info.Println(*f)
	request.Oldfilm = *f
	if request.NewFilm.Name != "" && len(request.NewFilm.Name) > 1 && len(request.NewFilm.Name) < 150 {
		s.Info.Println("editting name films")
		request.Oldfilm.Name = request.NewFilm.Name
	}
	if request.NewFilm.About != "" && len(request.NewFilm.About) < 1000 {
		s.Info.Println("editting about films")
		request.Oldfilm.About = request.NewFilm.About
	}
	if !request.NewFilm.ReleaseDate.IsZero() {
		s.Info.Println("editting realisedate films")
		request.Oldfilm.ReleaseDate = request.NewFilm.ReleaseDate
	}
	if request.NewFilm.Rating >= 0 && request.NewFilm.Rating <= 10 {
		s.Info.Println("editting ratring films")
		request.Oldfilm.Rating = request.NewFilm.Rating
	}
	e = s.DatabaseConnector.EditFilm(context.Background(), request.Oldfilm)
	if e != nil {
		s.Error.Println("error edit one film")
		return e
	}
	s.Info.Println("finish editing one filmrequest")
	return nil
}

func (s *Service) deltefilm(film entity.Film) error {
	s.Info.Println("start delete one film ")
	if film.Id != 0 {
		s.Error.Println("invalide film data ")
		return errors.New("invalide film data")

	}
	s.Info.Println(film)
	s.Info.Println("get one film ")
	f, err := s.DatabaseConnector.GetFilm(context.Background(), film)
	if err != nil {
		s.Error.Println("error getting one film ")
		return err
	}
	s.Info.Println("finished get one film ")
	s.Info.Println(*f)
	s.Info.Println("delete one film ")

	e := s.DatabaseConnector.DeleteFilm(context.Background(), *f)
	if e != nil {
		s.Error.Println("error delete one film ")
		return e
	}
	s.Info.Println("finished delete one film ")
	return nil
}

func (s *Service) AddConnectionBetweenActorAndFilm(data []byte) error {
	s.Info.Println("start add connection between film and actor")
	var con entity.RequestEditConnection
	err := json.Unmarshal(data, &con)
	if err != nil {
		s.Error.Println("error unmarshal RequestEditConnection in add connection")
		return err
	}
	tmp, err := s.getconnectionbyfilmandactor(con)
	if err != nil {
		s.Error.Println("error gets data from base for connection between film and actor")
		return err
	}
	con = *tmp
	if e := s.DatabaseConnector.AddActorFilmConnection(context.Background(), con.Actor.Id, con.Film.Id); e != nil {
		s.Error.Println("error add connection between film and actor")
		return e
	}
	s.Info.Println("finished add connection between film and actor")
	return nil
}

func (s *Service) DeleteConnectionBetweenActorAndFilm(data []byte) error {
	s.Info.Println("start delete connection between film and actor")
	var con entity.RequestEditConnection
	err := json.Unmarshal(data, &con)
	if err != nil {
		s.Error.Println("error unmarshal RequestEditConnection in delete connection")
		return err
	}
	tmp, err := s.getconnectionbyfilmandactor(con)
	if err != nil {
		s.Error.Println("error gets data from base for connection between film and actor")
		return err
	}
	con = *tmp
	if e := s.DatabaseConnector.DeleteActorFilmConnection(context.Background(), con.Actor.Id, con.Film.Id); e != nil {
		s.Error.Println("error delete connection between film and actor")
		return e
	}
	s.Info.Println("finished delete connection between film and actor")
	return nil

}

func (s *Service) getconnectionbyfilmandactor(req entity.RequestEditConnection) (*entity.RequestEditConnection, error) {
	s.Info.Println("start gettiing id film and actor from base")
	tmp, err := s.DatabaseConnector.GetActor(context.Background(), req.Actor)
	if err != nil {
		s.Error.Println("error get actor from base")
		return nil, err
	}
	req.Actor = *tmp
	tmp1, err := s.DatabaseConnector.GetFilm(context.Background(), req.Film)
	if err != nil {
		s.Error.Println("error get film from base")
		return nil, err
	}
	req.Film = *tmp1
	s.Info.Println("finish gettiing id film and actor from base")
	return &req, nil
}

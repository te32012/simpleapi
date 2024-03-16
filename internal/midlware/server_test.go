package midlware_test

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
	"vktestgo2024/internal/entity"
	"vktestgo2024/internal/midlware"
)

var list int = 0

type servicemock struct {
	hashadminb string
	hashuserb  string
}

func (s *servicemock) Login(login, password string) (string, error) {
	if login == "test" {
		return "test", nil
	}
	if login == "all" {
		return "", errors.New("произвольная ошибка")
	}
	return "", errors.New("нужного пользователя нет")
}
func (s *servicemock) CheckUserIsLoginedAndHasPermission(key string, operation int) bool {
	if key == s.hashadminb && operation == entity.AdminPermission {
		return true
	}
	if key == s.hashuserb && operation == entity.UserPermission {
		return true
	}
	return false
}
func (s *servicemock) AddActor(data []byte) error {
	fmt.Println(len(data))
	if len(data) == 1 && data[0] == 1 {
		return nil
	}
	return errors.New("тестовая ошибка addactor")
}
func (s *servicemock) EditActor(data []byte) error {
	fmt.Println(len(data))
	if len(data) == 1 && data[0] == 1 {
		return nil
	}
	return errors.New("тестовая ошибка editactor")
}
func (s *servicemock) DeleteActor(data []byte) error {
	fmt.Println(len(data))
	if len(data) == 1 && data[0] == 1 {
		return nil
	}
	return errors.New("тестовая ошибка deleteactor")
}
func (s *servicemock) GetListActors() ([]byte, error) {
	switch list {
	case 0:
		return nil, errors.New("тестовая ошибка getlistactors")
	case 1:
		return []byte{1}, nil
	default:
		return nil, errors.New("странно")
	}
}
func (s *servicemock) AddFilm(data []byte) error {
	fmt.Println(len(data))
	if len(data) == 1 && data[0] == 1 {
		return nil
	}
	return errors.New("тестовая ошибка editactor")
}
func (s *servicemock) EditFilm(data []byte) error {
	fmt.Println(len(data))
	if len(data) == 1 && data[0] == 1 {
		return nil
	}
	return errors.New("тестовая ошибка editactor")
}
func (s *servicemock) GetListFilms(keySort int, orderSort int) ([]byte, error) {
	switch list {
	case 0:
		return nil, errors.New("тестовая ошибка getlistfilms")
	case 1:
		return []byte{1}, nil
	default:
		return nil, errors.New("странно")
	}
}
func (s *servicemock) DeleteFilm(data []byte) error {
	fmt.Println(len(data))
	if len(data) == 1 && data[0] == 1 {
		return nil
	}
	return errors.New("тестовая ошибка editactor")
}
func (s *servicemock) FindInFilm(segment string) ([]byte, error) {
	if segment == "test" {
		return []byte{1}, nil
	}
	return nil, errors.New("тестовая ошибка findinfilm")
}

func (s *servicemock) AddConnectionBetweenActorAndFilm(data []byte) error {
	fmt.Println(len(data))
	if len(data) == 1 && data[0] == 1 {
		return nil
	}
	return errors.New("тестовая ошибка findinfilm")
}
func (s *servicemock) DeleteConnectionBetweenActorAndFilm(data []byte) error {
	fmt.Println(len(data))
	if len(data) == 1 && data[0] == 1 {
		return nil
	}
	return errors.New("тестовая ошибка findinfilm")
}

func TestRouter(t *testing.T) {
	ha := sha256.Sum256([]byte("admin" + time.Now().GoString()))
	tmp1 := fmt.Sprintf("%x", ha[:16])
	ha = sha256.Sum256([]byte("user" + time.Now().GoString()))
	tmp2 := fmt.Sprintf("%x", ha[:16])
	router := midlware.NewRouter("localhost", "2024", &servicemock{hashadminb: tmp1, hashuserb: tmp2})
	go router.Lisen()

	client := http.Client{}
	resp, err := client.Do(&http.Request{Method: http.MethodPost, URL: &url.URL{Host: "localhost:2024", Path: "/ping", Scheme: "http"}})
	if err != nil && resp.StatusCode != 200 {
		t.Fatal(err)
	}

	h := make(http.Header)
	h.Set("login", "test")
	h.Set("password", "password")
	resp, err = client.Do(&http.Request{Method: http.MethodPost, URL: &url.URL{Host: "localhost:2024", Path: "/login", Scheme: "http"}, Header: h})
	if err != nil {
		t.Fatal(err)
	}
	if !(resp.StatusCode == http.StatusOK && resp.Header.Get("session") == "test") {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}
	h.Set("login", "user")
	h.Set("password", tmp2)
	fmt.Println(h.Get("password"))
	resp, err = client.Do(&http.Request{Method: http.MethodPost, URL: &url.URL{Host: "localhost:2024", Path: "/login", Scheme: "http"}, Header: h})
	if err != nil {
		t.Fatal(err)
	}
	if !(resp.StatusCode == http.StatusBadRequest) {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}
	resp, err = client.Do(&http.Request{Method: http.MethodPost, URL: &url.URL{Host: "localhost:2024", Path: "/login", Scheme: "http"}, Header: h})
	if err != nil {
		t.Fatal(err)
	}
	h.Del("login")
	h.Set("password", "password")
	if !(resp.StatusCode == http.StatusBadRequest) {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ := http.NewRequest("GET", "http://localhost:2024/login", bytes.NewReader([]byte{1}))
	req.Header.Set("session", "user")
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/login", bytes.NewReader([]byte{1}))
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	h.Set("login", "all")
	h.Set("password", "all")

	req, _ = http.NewRequest("POST", "http://localhost:2024/login", bytes.NewReader([]byte{1}))
	req.Header = h
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/addActor", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp2)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusMethodNotAllowed && resp.Header.Get("session") == "user" {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/addActor", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp2)
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusForbidden && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/addActor", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusInternalServerError && resp.Header.Get("session") == tmp1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}
	req, _ = http.NewRequest("POST", "http://localhost:2024/addActor", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("GET", "http://localhost:2024/addActor", bytes.NewReader([]byte{0}))
	req.Header.Set("session", "admin")
	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusMethodNotAllowed && resp.Header.Get("session") == "admin" {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/editActor", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == "admin" {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/editActor", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusForbidden && resp.Header.Get("session") == "user" {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/editActor", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusInternalServerError && resp.Header.Get("session") == tmp1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("GET", "http://localhost:2024/editActor", bytes.NewReader([]byte{0}))
	req.Header.Set("session", "admin")

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusMethodNotAllowed && resp.Header.Get("session") == "admin" {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/deleteActor", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/deleteActor", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusInternalServerError && resp.Header.Get("session") == tmp1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}
	req, _ = http.NewRequest("GET", "http://localhost:2024/deleteActor", bytes.NewReader([]byte{0}))
	req.Header.Set("session", "admin")

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusMethodNotAllowed && resp.Header.Get("session") == "admin" {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/deleteActor", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusForbidden && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("GET", "http://localhost:2024/getActors", bytes.NewReader([]byte{0}))
	req.Header.Set("session", "user")

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusMethodNotAllowed && resp.Header.Get("session") == "user" {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	list = 0

	req, _ = http.NewRequest("POST", "http://localhost:2024/getActors", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)

	if resp.StatusCode != http.StatusInternalServerError && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	list = 1

	req, _ = http.NewRequest("POST", "http://localhost:2024/getActors", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)
	tmp, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("GET", "http://localhost:2024/addFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusMethodNotAllowed && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/addFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusForbidden && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/addFilm", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/addFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusInternalServerError && resp.Header.Get("session") == tmp1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("GET", "http://localhost:2024/editFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", "user")

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusMethodNotAllowed && resp.Header.Get("session") == "user" {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/editFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusForbidden && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/editFilm", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/editFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusInternalServerError && resp.Header.Get("session") == tmp1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/editFilm", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("GET", "http://localhost:2024/deleteFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", "user")

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusMethodNotAllowed && resp.Header.Get("session") == "user" {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/deleteFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusForbidden && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/deleteFilm", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/deleteFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusInternalServerError && resp.Header.Get("session") == tmp1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/deleteFilm", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("GET", "http://localhost:2024/findInFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusMethodNotAllowed && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/findInFilm", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusBadRequest && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/findInFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)
	req.Header.Set("segment", "test")

	resp, _ = client.Do(req)
	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}
	req, _ = http.NewRequest("POST", "http://localhost:2024/findInFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)
	req.Header.Set("segment", "e")

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusInternalServerError && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}
	req, _ = http.NewRequest("POST", "http://localhost:2024/findInFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", "1")
	req.Header.Set("segment", "e")

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusInternalServerError && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}
	req, _ = http.NewRequest("POST", "http://localhost:2024/findInFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", "1")
	req.Header.Set("segment", "e")

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusInternalServerError && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("GET", "http://localhost:2024/getListFilms", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)
	if resp.StatusCode != http.StatusMethodNotAllowed && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	list = 1

	req, _ = http.NewRequest("POST", "http://localhost:2024/getListFilms", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)
	req.Header.Set("keysort", "2")
	req.Header.Set("ordersort", "1")

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/getListFilms", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)
	req.Header.Set("keysort", "l")
	req.Header.Set("ordersort", "1")

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}
	req, _ = http.NewRequest("POST", "http://localhost:2024/getListFilms", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)
	req.Header.Set("keysort", "1")
	req.Header.Set("ordersort", "l")

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/getListFilms", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)
	req.Header.Set("keysort", "1")
	req.Header.Set("ordersort", "l")
	req, _ = http.NewRequest("POST", "http://localhost:2024/getListFilms", bytes.NewReader([]byte{0}))
	req.Header.Set("session", "1")
	req.Header.Set("keysort", "1")
	req.Header.Set("ordersort", "l")

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	list = 0

	req, _ = http.NewRequest("POST", "http://localhost:2024/getListFilms", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)
	req.Header.Set("keysort", "2")
	req.Header.Set("ordersort", "1")

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	list = 1

	req, _ = http.NewRequest("POST", "http://localhost:2024/get", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp2)
	req.Header.Set("keysort", "2")
	req.Header.Set("ordersort", "1")

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusNotFound && resp.Header.Get("session") == tmp2 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/addConnectionBetweenActorAndFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/addConnectionBetweenActorAndFilm", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("GET", "http://localhost:2024/addConnectionBetweenActorAndFilm", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("GET", "http://localhost:2024/addConnectionBetweenActorAndFilm", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/deleteConnectionBetweenActorAndFilm", bytes.NewReader([]byte{0}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/deleteConnectionBetweenActorAndFilm", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp1)

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("POST", "http://localhost:2024/deleteConnectionBetweenActorAndFilm", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

	req, _ = http.NewRequest("GET", "http://localhost:2024/deleteConnectionBetweenActorAndFilm", bytes.NewReader([]byte{1}))
	req.Header.Set("session", tmp2)

	resp, _ = client.Do(req)

	tmp, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.Header.Get("session") == tmp2 && len(tmp) == 1 && tmp[0] == 1 {
		t.Fatal(resp.StatusCode, resp.Header.Get("session"))
	}

}

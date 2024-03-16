package midlware

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"vktestgo2024/internal/entity"
	"vktestgo2024/internal/service"
)

type Router struct {
	Server  *http.Server
	Info    *log.Logger
	Error   *log.Logger
	Service service.ServiceInterface
}

func NewRouter(host string, port string, service service.ServiceInterface) *Router {
	r := &Router{}
	r.Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	r.Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.notfund)
	r.Info.Println("added / midlware")

	mux.HandleFunc("/ping", r.ping)
	r.Info.Println("added /ping midlware")

	mux.HandleFunc("/login", r.login)
	r.Info.Println("added /login midlware")

	mux.HandleFunc("/addActor", r.addActor)
	r.Info.Println("added /addActor midlware")
	mux.HandleFunc("/editActor", r.editActor)
	r.Info.Println("added /editActor midlware")
	mux.HandleFunc("/deleteActor", r.deleteActor)
	r.Info.Println("added /deleteActor midlware")
	mux.HandleFunc("/getActors", r.getActors)
	r.Info.Println("added /getActors midlware")

	mux.HandleFunc("/addFilm", r.addFilm)
	r.Info.Println("added /addFilm midlware")
	mux.HandleFunc("/editFilm", r.editFilm)
	r.Info.Println("added /editFilm midlware")
	mux.HandleFunc("/deleteFilm", r.deleteFilm)
	r.Info.Println("added /deleteFilm midlware")
	mux.HandleFunc("/getListFilms", r.getListFilms)
	r.Info.Println("added /getListFilms midlware")
	mux.HandleFunc("/findInFilm", r.findInFilm)
	r.Info.Println("added /findInFilm midlware")
	r.Server = &http.Server{
		Addr:    host + ":" + port,
		Handler: mux,
	}
	r.Service = service
	return r
}
func (r *Router) ping(response http.ResponseWriter, request *http.Request) {
	r.Info.Println("request on midlware " + request.RequestURI + " using method " + request.Method)
	response.WriteHeader(http.StatusOK)
	response.Write([]byte(fmt.Sprintf("<h1>200 OK</h1><p>%s</p>", "сервер успешно работает!!!")))
	response.Header().Set("Content-Type", "text/html")
}

func (r *Router) notfund(response http.ResponseWriter, request *http.Request) {
	r.Info.Println("request on midlware " + request.RequestURI + " using method " + request.Method)
	r.Error.Println("server not listen this URI aders ")
	response.WriteHeader(http.StatusNotFound)
	response.Write([]byte(fmt.Sprintf("<h1>404 Not Found</h1><p>%s</p>", "endpoint не существует")))
	response.Header().Set("Content-Type", "text/html")
}
func (r *Router) login(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	if !r.checkrequestmethond(response, request) {
		return
	}

	login := request.Header.Get("login")
	password := request.Header.Get("password")
	if login == "" || password == "" {
		r.Error.Println("нет login или password")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(fmt.Sprintf("<h1>400 Bad Request</h1><p>%s</p>", "пустой логин или пароль")))
		response.Header().Set("Content-Type", "text/html")
		return
	}
	session, err := r.Service.Login(login, password)

	if err != nil && err.Error() == "нужного пользователя нет" {
		r.Error.Println("ошибочный логин или пароль")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(fmt.Sprintf("<h1>400 Bad Request</h1><p>%s</p>", err)))
		response.Header().Set("Content-Type", "text/html")
		return
	}
	if err != nil {
		r.Error.Println(err)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("<h1>500 Internal Server Error</h1><p>%s</p>", err)))
		response.Header().Set("Content-Type", "text/html")
		return
	}
	r.Info.Println("user " + login + " session " + session + " was logined")
	response.Header().Set("session", session)
}

func (r *Router) addActor(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	if !r.checkrequestmethond(response, request) {
		return
	}
	if !r.checkrequestpermission(response, request, entity.AdminPermission) {
		return
	}
	data, ok := r.getdatafromrequest(response, request)
	if !ok {
		return
	}
	e := r.Service.AddActor(data)
	if e != nil {
		r.Error.Println("incorrect response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
		r.Error.Println(e)
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("<h1>500 Internal Server Error</h1><p>%s</p>", e)))
		return
	}
	r.Info.Println("correct response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
}

func (r *Router) editActor(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	if !r.checkrequestmethond(response, request) {
		return
	}

	if !r.checkrequestpermission(response, request, entity.AdminPermission) {
		return
	}
	data, ok := r.getdatafromrequest(response, request)
	if !ok {
		return
	}
	e := r.Service.EditActor(data)
	if e != nil {
		r.Error.Println("incorrect response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
		r.Error.Println(e)
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("<h1>500 Internal Server Error</h1><p>%s</p>", e)))
		return
	}
	r.Info.Println("correct response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
}

func (r *Router) deleteActor(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	if !r.checkrequestmethond(response, request) {
		return
	}

	if !r.checkrequestpermission(response, request, entity.AdminPermission) {
		return
	}
	data, ok := r.getdatafromrequest(response, request)
	if !ok {
		return
	}
	e := r.Service.DeleteActor(data)
	if e != nil {
		r.Error.Println("incorrect response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
		r.Error.Println(e)
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("<h1>500 Internal Server Error</h1><p>%s</p>", e)))
		return
	}
	r.Info.Println("correct response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
}

func (r *Router) getActors(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	if !r.checkrequestmethond(response, request) {
		return
	}

	if !r.checkrequestpermission(response, request, entity.UserPermission) {
		return
	}
	data, e := r.Service.GetListActors()
	if e != nil {
		r.Error.Println("incorrect response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
		r.Error.Println(e)
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("<h1>500 Internal Server Error</h1><p>%s</p>", e)))
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.Write(data)
	r.Info.Println("correct response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
}

func (r *Router) addFilm(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	if !r.checkrequestmethond(response, request) {
		return
	}
	if !r.checkrequestpermission(response, request, entity.AdminPermission) {
		return
	}
	data, ok := r.getdatafromrequest(response, request)
	if !ok {
		return
	}

	e := r.Service.AddFilm(data)
	if e != nil {
		r.Error.Println("incorrect response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
		r.Error.Println(e)
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("<h1>500 Internal Server Error</h1><p>%s</p>", e)))
		return
	}
	r.Info.Println("correct response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
}

func (r *Router) editFilm(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	if !r.checkrequestmethond(response, request) {
		return
	}

	if !r.checkrequestpermission(response, request, entity.AdminPermission) {
		return
	}
	data, ok := r.getdatafromrequest(response, request)
	if !ok {
		return
	}

	e := r.Service.EditFilm(data)
	if e != nil {
		r.Error.Println("incorrect response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
		r.Error.Println(e)
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("<h1>500 Internal Server Error</h1><p>%s</p>", e)))
		return
	}
	r.Info.Println("correct response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
}

func (r *Router) deleteFilm(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	if !r.checkrequestmethond(response, request) {
		return
	}
	if !r.checkrequestpermission(response, request, entity.AdminPermission) {
		return
	}
	data, ok := r.getdatafromrequest(response, request)
	if !ok {
		return
	}

	e := r.Service.DeleteFilm(data)
	if e != nil {
		r.Error.Println("incorrect response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
		r.Error.Println(e)
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("<h1>500 Internal Server Error</h1><p>%s</p>", e)))
		return
	}
	r.Info.Println("correct response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
}

func (r *Router) getListFilms(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	if !r.checkrequestmethond(response, request) {
		return
	}
	if !r.checkrequestpermission(response, request, entity.UserPermission) {
		return
	}
	var keysort int = 3
	var ordersort int = -1
	ks := request.Header.Get("keysort")
	if ks != "" {
		var err error
		keysort, err = strconv.Atoi(ks)
		if err != nil {
			r.Error.Println("incorrect response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
			response.Header().Set("Content-Type", "text/html")
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(fmt.Sprintf("<h1>500 Internal Server Error</h1><p>%s</p>", "некорректный keysort")))
			return
		}
	}
	os := request.Header.Get("ordersort")
	if os != "" {
		var err error
		ordersort, err = strconv.Atoi(os)
		if err != nil {
			r.Error.Println("incorrect response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
			response.Header().Set("Content-Type", "text/html")
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte(fmt.Sprintf("<h1>500 Internal Server Error</h1><p>%s</p>", "некорректный ordersort")))
			return
		}
	}

	ans, e := r.Service.GetListFilms(keysort, ordersort)
	if e != nil {
		r.Error.Println("incorrect response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
		r.Error.Println(e)
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("<h1>500 Internal Server Error</h1><p>%s</p>", e)))
		return
	}
	response.Write(ans)
	response.Header().Set("Content-Type", "application/json")
	r.Info.Println("correct response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
}

func (r *Router) findInFilm(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	if !r.checkrequestmethond(response, request) {
		return
	}
	if !r.checkrequestpermission(response, request, entity.UserPermission) {
		return
	}

	if request.Header.Get("segment") == "" {
		r.Error.Println("incorrect response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(fmt.Sprintf("<h1>500 Internal Server Error</h1><p>%s</p>", "segment не задан")))
		return
	}

	ans, err := r.Service.FindInFilm(request.Header.Get("segment"))
	if err != nil {
		r.Error.Println("incorrect response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
		r.Error.Println(err)
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("<h1>500 Internal Server Error</h1><p>%s</p>", err)))
		return
	}
	response.Write(ans)
	response.Header().Set("Content-Type", "application/json")
	r.Info.Println("correct response for uri " + request.RequestURI + " and session " + request.Header.Get("session"))
}

func (r *Router) Lisen() {
	r.Info.Println("starting server on host:port " + r.Server.Addr)
	log.Fatal(r.Server.ListenAndServe())
	r.Info.Println("finised server")
}

func (r *Router) checkrequestmethond(response http.ResponseWriter, request *http.Request) bool {
	r.Info.Println("request on midlware " + request.RequestURI + " using method " + request.Method)

	response.Header().Set("session", request.Header.Get("session"))

	if request.Method != http.MethodPost {
		r.Error.Println("method not allowed")
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte(fmt.Sprintf("<h1>405 Method Not Allowed</h1><p>%s</p>", "method not allowed")))
		return false
	}
	return true
}

func (r *Router) checkrequestpermission(response http.ResponseWriter, request *http.Request, permisstion int) bool {
	r.Info.Println("check permission for uri " + request.RequestURI + " for session " + request.Header.Get("session"))

	if len(request.Header.Get("session")) != 32 {
		r.Error.Println("сессия не задана корректно")
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusUnauthorized)
		response.Write([]byte(fmt.Sprintf("<h1>401 StatusUnauthorized</h1><p>%s</p>", "method not allowed")))
		return false
	}
	ok := r.Service.CheckUserIsLoginedAndHasPermission(request.Header.Get("session"), permisstion)
	if !ok {
		r.Error.Println("доступ запрещен")
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusForbidden)
		response.Write([]byte(fmt.Sprintf("<h1>403 StatusForbidden</h1><p>%s</p>", "доступ запрещен")))
		return false
	}
	r.Info.Println("checked permission for uri " + request.RequestURI + " for session " + request.Header.Get("session"))
	return true
}

func (r *Router) getdatafromrequest(response http.ResponseWriter, request *http.Request) ([]byte, bool) {
	r.Info.Println("save data for uri " + request.RequestURI + " for session " + request.Header.Get("session"))

	data, e := io.ReadAll(request.Body)
	if e != nil {
		r.Error.Println(e)
		response.Header().Set("Content-Type", "text/html")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(fmt.Sprintf("<h1>500 StatusInternalServerError</h1><p>%s</p>", e)))
		return nil, false
	}
	r.Info.Println("saved data for uri " + request.RequestURI + " for session " + request.Header.Get("session"))
	return data, true
}

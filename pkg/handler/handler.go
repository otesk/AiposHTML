package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Controller interface {
	Root(url string) ([]byte, error)
}

type Handler struct {
	C *http.Client
}

func BuildRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.welcome)
	mux.HandleFunc("/user/new", h.userPage)
	mux.HandleFunc("/user/create", h.createUser)
	mux.HandleFunc("/user/find", h.user)
	mux.HandleFunc("/user/fund", h.fund)
	mux.HandleFunc("/user/delete", h.deleteUser)
	mux.HandleFunc("/tournament/new", h.tournamentPage)
	mux.HandleFunc("/tournament/create", h.createTournament)
	mux.HandleFunc("/tournament/find", h.tournament)
	mux.HandleFunc("/tournament/finish", h.finishTournament)
	mux.HandleFunc("/tournament/cancel", h.cancelTournament)
	mux.HandleFunc("/tournament/join", h.join)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	return mux
}

func (h Handler) welcome(w http.ResponseWriter, r *http.Request) {

	temp, err := template.ParseFiles("static/templates/greet.html")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = temp.Execute(w, nil)

}

func (h Handler) userPage(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("static/templates/newUser.html")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = temp.Execute(w, nil)
}

func (h Handler) tournamentPage(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("static/templates/newTournament.html")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = temp.Execute(w, nil)
}

func (h Handler) createUser(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("username")

	data, err := json.Marshal(map[string]string{"name": name})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.C.Post("http://users-service:8080/user", "Application/JSON", bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	var user user

	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte("Great news! User with an id of " + user.ID + " Has been created!"))
	if err != nil {
		log.Println(err)
	}
}

func (h Handler) user(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("user")

	resp, err := h.C.Get("http://users-service:8080/user/" + id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	var user user
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		log.Println(err)
		log.Println(user)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	temp, err := template.ParseFiles("static/templates/user.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = temp.Execute(w, &user)
	if err != nil {
		log.Println(err)
	}

}

func (h Handler) fund(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("user")
	p := r.URL.Query().Get("points")
	points, err := strconv.Atoi(p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(map[string]int{"points": points})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.C.
		Post("http://users-service:8080/user/"+id+"/fund", "Application/JSON", bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	fmt.Println(resp.Status)

	_, err = w.Write([]byte("Funded " + p + " points to " + id))
	if err != nil {
		log.Println(err)
	}

}

func (h Handler) take(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("user")
	p := r.URL.Query().Get("points")
	points, err := strconv.Atoi(p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(map[string]int{"points": points})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.C.
		Post("http://users-service:8080/user/"+id+"/take", "Application/JSON", bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	fmt.Println(resp.Status)

	_, err = w.Write([]byte("Funded " + p + " points to " + id))
	if err != nil {
		log.Println(err)
	}

}

func (h Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("user")

	req, err := http.NewRequest(http.MethodDelete, "http://users-service:8080/user"+id, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.C.Do(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	_, err = w.Write([]byte("Great news! User with an id of " + id + " Has been deleted"))
	if err != nil {
		log.Println(err)
	}
}

func (h Handler) createTournament(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	d := r.URL.Query().Get("deposit")
	deposit, err := strconv.Atoi(d)
	if err != nil {
		log.Println(err)
		return
	}

	data, err := json.Marshal(map[string]interface{}{"name": name, "deposit": deposit})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.C.Post("http://tournament-service:8081/tournament", "Application/JSON", bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	var tournament tournamentData

	err = json.NewDecoder(resp.Body).Decode(&tournament)
	if err != nil {
		log.Println(err)
		log.Println(tournament)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte("Great news! Tournament with an id of " + tournament.ID + " Has been created!"))
	if err != nil {
		log.Println(err)
	}
}

func (h Handler) tournament(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("tournament")

	resp, err := h.C.Get("http://tournament-service:8081/tournament/" + id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	var tournamentData tournamentData
	err = json.NewDecoder(resp.Body).Decode(&tournamentData)
	if err != nil {
		log.Println(err)
		log.Println(tournamentData)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	temp, err := template.ParseFiles("static/templates/tournament.html")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = temp.Execute(w, &tournamentData)
	if err != nil {
		log.Println(err)
	}

}

func (h Handler) cancelTournament(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("tournament")

	req, err := http.NewRequest(http.MethodDelete, "http://tournament-service:8081/tournament"+id, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.C.Do(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	_, err = w.Write([]byte("Great news! Tournament with an id of " + id + " Has been cancelled"))
	if err != nil {
		log.Println(err)
	}
}

func (h Handler) finishTournament(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("tournament")

	req, err := http.NewRequest(http.MethodPost, "http://tournament-service:8081/tournament"+id+"/finish", nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.C.Do(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	_, err = w.Write([]byte("Great news! Tournament with an id of " + id + " Has been finished"))
	if err != nil {
		log.Println(err)
	}
}

func (h Handler) join(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("tournament")
	user := r.URL.Query().Get("user")

	data, err := json.Marshal(map[string]string{"id": user})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.C.
		Post("http://tournament-service:8081/tournament/"+id+"/join", "Application/JSON", bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	fmt.Println(resp.Status)

	_, err = w.Write([]byte("User " + user + " joined tournament " + id))
	if err != nil {
		log.Println(err)
	}

}

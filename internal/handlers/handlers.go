package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/redblue-blur/bookings/internal/config"
	"github.com/redblue-blur/bookings/internal/forms"
	"github.com/redblue-blur/bookings/internal/models"
	"github.com/redblue-blur/bookings/internal/render"
)

//Repo is the repository used by handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

//NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

//newHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home is the home page handeler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
}

//About is the about page handeler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "hello again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

//Reservation is the Reservation form handeler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

//PostReservation handels trhe posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}
	form := forms.New(r.PostForm)
	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}
}

//Generals is the Room page handeler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.html", &models.TemplateData{})
}

//Majors is the Room page handeler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.html", &models.TemplateData{})
}

//Availability is the search availability page handeler
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.html", &models.TemplateData{})
}

//PostAvailability is the search availability page handeler
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s end date is %s", start, end)))
	// render.RenderTemplate(w, r, "search-availability.page.html", &models.TemplateData{})
}

type JsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

//AvailabilityJSON is the requests for availability and sends Json response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := JsonResponse{
		Ok:      true,
		Message: "available",
	}
	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

//Contact is the Contact page handeler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})
}

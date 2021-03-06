package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redblue-blur/bookings/internal/config"
	"github.com/redblue-blur/bookings/internal/driver"
	"github.com/redblue-blur/bookings/internal/forms"
	"github.com/redblue-blur/bookings/internal/helpers"
	"github.com/redblue-blur/bookings/internal/models"
	"github.com/redblue-blur/bookings/internal/render"
	"github.com/redblue-blur/bookings/internal/repository"
	"github.com/redblue-blur/bookings/internal/repository/dbrepo"
)

//Repo is the repository used by handlers
var Repo *Repository

//Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

//newHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home is the home page handeler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.html", &models.TemplateData{})
}

//About is the about page handeler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//send data to template
	render.Template(w, r, "about.page.html", &models.TemplateData{})
}

//Reservation is the Reservation form handeler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "reservation.page.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

//PostReservation handels trhe posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
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
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "reservation.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

//ReservationSummary is used to display reservation details
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.html", &models.TemplateData{
		Data: data,
	})
}

//Generals is the Room page handeler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.html", &models.TemplateData{})
}

//Majors is the Room page handeler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.html", &models.TemplateData{})
}

//Availability is the search availability page handeler
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.html", &models.TemplateData{})
}

//PostAvailability is the search availability page handeler
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is %s end date is %s", start, end)))
	// render.Template(w, r, "search-availability.page.html", &models.TemplateData{})
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
		helpers.ServerError(w, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

//Contact is the Contact page handeler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

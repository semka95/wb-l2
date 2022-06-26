package api

import (
	"calendar/event"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type API struct {
	eventStore event.EventRepository
	logger     *zap.Logger
}

func NewAPI(repository event.EventRepository, logger *zap.Logger) API {
	return API{
		eventStore: repository,
		logger:     logger,
	}
}

func (a *API) NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", a.Create)
	mux.HandleFunc("/update_event", a.Update)
	mux.HandleFunc("/delete_event", a.Delete)
	mux.HandleFunc("/events_for_day", a.Get)
	mux.HandleFunc("/events_for_week", a.Get)
	mux.HandleFunc("/events_for_month", a.Get)

	return mux
}

func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("bad method: %s", r.Method), "method should be post")
		return
	}

	err := r.ParseForm()
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse form")
		return
	}

	uid := r.FormValue("user_id")
	user_id, err := strconv.Atoi(uid)
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse user_id")
		return
	}

	date := r.FormValue("date")
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse date, use RFC3339 format")
		return
	}

	title := r.FormValue("title")
	if title == "" {
		SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("empty title"), "no title provided")
		return
	}

	e := event.Event{
		Title: title,
		Date:  t,
	}

	err = a.eventStore.Create(uint64(user_id), e)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't create event")
		return
	}

	SendNoContent(w, r)
}

func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("bad method: %s", r.Method), "method should be put")
		return
	}

	err := r.ParseForm()
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse form")
		return
	}

	uid := r.FormValue("user_id")
	user_id, err := strconv.Atoi(uid)
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse user_id")
		return
	}

	date := r.FormValue("date")
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse date, use RFC3339 format")
		return
	}

	title := r.FormValue("title")
	if title == "" {
		SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("empty title"), "no title provided")
		return
	}

	e := event.Event{
		Title: title,
		Date:  t,
	}

	err = a.eventStore.Update(uint64(user_id), e)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't update event")
		return
	}

	SendNoContent(w, r)
}

func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("bad method: %s", r.Method), "method should be delete")
		return
	}

	err := r.ParseForm()
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse form")
		return
	}

	uid := r.FormValue("user_id")
	user_id, err := strconv.Atoi(uid)
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse user_id")
		return
	}

	eid := r.FormValue("id")
	event_id, err := strconv.Atoi(eid)
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse id")
		return
	}

	err = a.eventStore.Delete(uint64(user_id), uint64(event_id))
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't delete event")
		return
	}

	SendNoContent(w, r)
}

func (a *API) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("bad method: %s", r.Method), "method should be get")
		return
	}

	uid := r.URL.Query().Get("user_id")
	user_id, err := strconv.Atoi(uid)
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse user_id")
		return
	}

	date := r.URL.Query().Get("date")
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse date, use RFC3339 format")
		return
	}

	events := make([]event.Event, 0)
	switch r.URL.Path {
	case "/events_for_day":
		events, err = a.eventStore.GetForDay(uint64(user_id), t)
	case "/events_for_week":
		events, err = a.eventStore.GetForWeek(uint64(user_id), t)
	case "/events_for_month":
		events, err = a.eventStore.GetForMonth(uint64(user_id), t)
	}

	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get events")
		return
	}

	SendJSON(w, r, http.StatusOK, events)
}
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"task11/repository"
	"task11/storage"
	"time"
)

type Handler struct {
	db *storage.EventStorage
}

func NewHandler() *Handler {
	return &Handler{
		db: storage.NewEventStorage(),
	}
}

type ResultResponse struct {
	Result []repository.Event `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event repository.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		SendErrorResponse(w, err, http.StatusServiceUnavailable)
		return
	}
	err = h.db.Insert(event)
	if err != nil {
		SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	SendResultResponse(w, []repository.Event{event})
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var event repository.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		SendErrorResponse(w, err, http.StatusServiceUnavailable)
		return
	}
	err = h.db.Update(event)
	if err != nil {
		SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	SendResultResponse(w, []repository.Event{event})
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var event repository.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		SendErrorResponse(w, err, http.StatusServiceUnavailable)
		return
	}
	h.db.Delete(event)
}

func (h *Handler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	userId, date, err := ParseQuery(r)
	if err != nil {
		SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	events := h.db.FindForDay(userId, date)
	SendResultResponse(w, events)
}

func (h *Handler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	userId, date, err := ParseQuery(r)
	if err != nil {
		SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	events := h.db.FindForWeek(userId, date)
	SendResultResponse(w, events)
}

func (h *Handler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	userId, date, err := ParseQuery(r)
	if err != nil {
		SendErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	events := h.db.FindForMonth(userId, date)
	SendResultResponse(w, events)
}

func ParseDate(date string) (time.Time, error) {
	eventDate, err := time.Parse("2006-01-02T15:04:00Z", date)
	if err != nil {
		eventDate, err = time.Parse("2006-01-02T15:04", date)
		if err != nil {
			eventDate, err = time.Parse("2006-01-02", date)
			if err != nil {
				return time.Time{}, fmt.Errorf("%v", err)
			}
		}
	}
	return eventDate, nil
}

func ParseQuery(r *http.Request) (int, time.Time, error) {
	userId, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil || userId < 1 {
		if userId < 1 {
			return -1, time.Time{}, errors.New("invalid user ID")
		}
		return -1, time.Time{}, err
	}
	date, err := ParseDate(r.URL.Query().Get("date"))
	if err != nil {
		return -1, time.Time{}, err
	}
	return userId, date, nil
}

func SendResultResponse(w http.ResponseWriter, events []repository.Event) {
	jsonData, _ := json.MarshalIndent(ResultResponse{Result: events}, "", " ")
	_, err := w.Write(jsonData)
	if err != nil {
		SendErrorResponse(w, err, http.StatusServiceUnavailable)
	}
}

func SendErrorResponse(w http.ResponseWriter, err error, status int) {
	jsonData, _ := json.MarshalIndent(ErrorResponse{Error: err.Error()}, "", " ")
	http.Error(w, string(jsonData), status)
}

func (h *Handler) InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/create_event", h.CreateEvent)
	mux.HandleFunc("/update_event", h.UpdateEvent)
	mux.HandleFunc("/delete_event", h.DeleteEvent)
	mux.HandleFunc("/events_for_day", h.GetEventsForDay)
	mux.HandleFunc("/events_for_week", h.GetEventsForWeek)
	mux.HandleFunc("/events_for_month", h.GetEventsForMonth)
}

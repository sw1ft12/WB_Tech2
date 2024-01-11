package storage

import (
	"errors"
	"task11/repository"
	"time"
)

type EventStorage struct {
	data map[int]map[int]repository.Event
}

func NewEventStorage() *EventStorage {
	return &EventStorage{
		data: make(map[int]map[int]repository.Event),
	}
}

func (s *EventStorage) Insert(event repository.Event) error {
	if s.data[event.UserId] == nil {
		s.data[event.UserId] = make(map[int]repository.Event)
	} else if _, ok := s.data[event.UserId][event.EventId]; ok {
		return errors.New("event already exists")
	}
	s.data[event.UserId][event.EventId] = event
	return nil
}

func (s *EventStorage) Update(event repository.Event) error {
	if s.data[event.UserId] == nil {
		s.data[event.UserId] = make(map[int]repository.Event)
	} else if _, ok := s.data[event.UserId][event.EventId]; !ok {
		return errors.New("event does not exist")
	}
	s.data[event.UserId][event.EventId] = event
	return nil
}

func (s *EventStorage) Delete(event repository.Event) {
	if s.data[event.UserId] == nil {
		s.data[event.UserId] = make(map[int]repository.Event)
	}
	delete(s.data[event.UserId], event.EventId)
}

func (s *EventStorage) FindForDay(userId int, date time.Time) []repository.Event {
	var result []repository.Event
	for _, event := range s.data[userId] {
		if event.Date.Compare(date) == 0 {
			result = append(result, event)
		}
	}
	return result
}

func (s *EventStorage) FindForWeek(userId int, date time.Time) []repository.Event {
	var result []repository.Event
	year, week := date.ISOWeek()
	for _, event := range s.data[userId] {
		eventYear, eventWeek := event.Date.ISOWeek()
		if year == eventYear && week == eventWeek {
			result = append(result, event)
		}
	}
	return result
}

func (s *EventStorage) FindForMonth(userId int, date time.Time) []repository.Event {
	year, month, _ := date.Date()
	var result []repository.Event
	for _, event := range s.data[userId] {
		eventYear, eventMonth, _ := event.Date.Date()
		if year == eventYear && eventMonth == month {
			result = append(result, event)
		}
	}
	return result
}

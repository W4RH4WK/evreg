// evreg is simple event registry allowing people to register for attendee
// limited slots.
// Copyright (C) 2021  Alexander Hirsch
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package evreg

import (
	"errors"
)

type EventRegister struct {
	Events []Event
}

func (e *EventRegister) FindEvent(token string) (*Event, error) {
	for _, event := range e.Events {
		if token == event.RegisterToken {
			return &event, nil
		}
	}
	return nil, errors.New("no such event")
}

func (e *EventRegister) FindEventAdmin(token string) (*Event, error) {
	for _, event := range e.Events {
		if token == event.AdminToken {
			return &event, nil
		}
	}
	return nil, errors.New("no such event")
}

func (e *EventRegister) Store(filepath string) error {
	return marshalToFile(filepath, e)
}

func LoadEventRegister(filepath string) (e EventRegister, err error) {
	err = unmarshalFromFile(filepath, &e)
	return
}

type Event struct {
	Name          string
	RegisterToken string
	AdminToken    string
	Slots         []Slot
}

func (e *Event) IsRegistered(name string) bool {
	for _, slot := range e.Slots {
		if slot.IsRegistered(name) {
			return true
		}
	}
	return false
}

type Slot struct {
	Name      string
	Limit     int
	Attendees []string
}

func (s *Slot) IsFull() bool {
	return len(s.Attendees) >= s.Limit
}

func (s *Slot) IsRegistered(name string) bool {
	for _, attendee := range s.Attendees {
		if name == attendee {
			return true
		}
	}
	return false
}

func (s *Slot) AddAttendee(name string) error {
	if s.IsFull() {
		return errors.New("slot already full")
	}

	if s.IsRegistered(name) {
		return errors.New("already registered")
	}

	s.Attendees = append(s.Attendees, name)
	return nil
}

func (s *Slot) RemoveAttendee(name string) error {
	for i, attendee := range s.Attendees {
		if attendee == name {
			s.Attendees = append(s.Attendees[:i], s.Attendees[i+1:]...)
			return nil
		}
	}
	return errors.New(name + " not registered")
}

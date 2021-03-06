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

package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/W4RH4WK/evreg/pkg/evreg"

	"github.com/gorilla/mux"
)

const (
	registerFilepath = "register.json"
)

var (
	register     evreg.EventRegister
	registerLock sync.RWMutex
	nameRegex    = regexp.MustCompile(`[0-9]+`)
	templates    = make(map[string]*template.Template)
)

//////////////////////////////////////////////////////////////////////////

func eventFromRequest(r *http.Request) (*evreg.Event, error) {
	return register.FindEvent(mux.Vars(r)["token"])
}

func eventFromRequestAdmin(r *http.Request) (*evreg.Event, error) {
	return register.FindEventAdmin(mux.Vars(r)["token"])
}

func slotFromRequestWithEvent(r *http.Request, event *evreg.Event) (*evreg.Slot, *evreg.Event, error) {
	vars := mux.Vars(r)
	slotIndex, err := strconv.Atoi(vars["slot"])
	if err != nil {
		return nil, nil, err
	}

	if slotIndex >= len(event.Slots) {
		return nil, nil, errors.New("invalid slot")
	}

	return &event.Slots[slotIndex], event, nil
}

func slotFromRequest(r *http.Request) (*evreg.Slot, *evreg.Event, error) {
	event, err := eventFromRequest(r)
	if err != nil {
		return nil, nil, err
	}
	return slotFromRequestWithEvent(r, event)
}

func slotFromRequestAdmin(r *http.Request) (*evreg.Slot, *evreg.Event, error) {
	event, err := eventFromRequestAdmin(r)
	if err != nil {
		return nil, nil, err
	}
	return slotFromRequestWithEvent(r, event)
}

func viewEvent(w http.ResponseWriter, r *http.Request) {
	registerLock.RLock()
	defer registerLock.RUnlock()

	event, err := eventFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = templates["event.html"].Execute(w, event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func viewEventAdmin(w http.ResponseWriter, r *http.Request) {
	registerLock.RLock()
	defer registerLock.RUnlock()

	event, err := eventFromRequestAdmin(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = templates["event_admin.html"].Execute(w, event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func showRegisterForm(w http.ResponseWriter, r *http.Request) {
	registerLock.RLock()
	defer registerLock.RUnlock()

	slot, _, err := slotFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = templates["register.html"].Execute(w, slot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addAttendeeToSlot(w http.ResponseWriter, r *http.Request) {
	registerLock.Lock()
	defer registerLock.Unlock()

	slot, event, err := slotFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	name := r.FormValue("name")
	if !nameRegex.MatchString(name) {
		http.Error(w, "invalid name", http.StatusBadRequest)
		return
	}

	if event.IsRegistered(name) {
		http.Error(w, "already registered", http.StatusForbidden)
		return
	}

	err = slot.AddAttendee(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = register.Store(registerFilepath)
	if err != nil {
		log.Println("register store:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		EventToken string
		SlotName   string
	}{
		event.RegisterToken,
		slot.Name,
	}

	err = templates["register_successful.html"].Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func removeAttendeeFromSlot(w http.ResponseWriter, r *http.Request) {
	registerLock.Lock()
	defer registerLock.Unlock()

	slot, _, err := slotFromRequestAdmin(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	attendee := mux.Vars(r)["attendee"]

	err = slot.RemoveAttendee(attendee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = register.Store(registerFilepath)
	if err != nil {
		log.Println("register store:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templates["delete_successful.html"].Execute(w, attendee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//////////////////////////////////////////////////////////////////////////

func loadTemplates() error {
	files, err := filepath.Glob("res/tpl/*.tpl")
	if err != nil {
		return err
	}

	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
		templates[name], err = template.ParseFiles(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	log.Println("EvReg version 0.1")

	var err error

	if err = loadTemplates(); err != nil {
		log.Fatalln(err.Error())
	}

	if register, err = evreg.LoadEventRegister(registerFilepath); err != nil {
		log.Fatalln(err.Error())
	}

	addr := ":8080"
	if len(os.Args) >= 2 {
		addr = os.Args[1]
	}

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/event/{token}", viewEvent).Methods(http.MethodGet)
	router.HandleFunc("/event/{token}/register/{slot}", showRegisterForm).Methods(http.MethodGet)
	router.HandleFunc("/event/{token}/register/{slot}", addAttendeeToSlot).Methods(http.MethodPost)
	router.HandleFunc("/event/{token}/admin", viewEventAdmin)
	router.HandleFunc("/event/{token}/admin/{slot}/delete/{attendee}", removeAttendeeFromSlot)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./res/web")))

	http.Handle("/", router)
	log.Printf("Listening on %s\n", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

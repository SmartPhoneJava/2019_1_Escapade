package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	ih "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"github.com/gorilla/mux"
)

/*
Handler contains all API operations
DB - database, where api work with information
Cookie - cookie settings, more in structure config.Cookie
Clients - grps.Clients, need to connect to Auth server
Auth - auth settinds, more in structure config.Auth
*/
type Handler struct {
	Cookie     config.Cookie
	Clients    *clients.Clients
	AuthClient config.AuthClient
	Auth       config.Auth
}

// Init configuration fields
func (h *Handler) Init(c *config.Configuration) {
	h.Cookie = c.Cookie
	h.AuthClient = c.AuthClient
	h.Auth = c.Auth
}

func (h *Handler) getUserID(r *http.Request) (int, error) {
	return ih.IntFromPath(r, "id", 1, re.ErrorInvalidUserID())
}

func (h *Handler) getPage(r *http.Request) int {

	page, _ := ih.IntFromPath(r, "page", 1, nil)
	return page
}

func (h *Handler) getPerPage(r *http.Request) int {

	page, _ := ih.IntFromPath(r, "per_page", 100, nil)
	return page
}

func (h *Handler) getDifficult(r *http.Request) int {

	diff, _ := ih.IntFromPath(r, "difficult", 0, nil)
	if diff > 3 {
		diff = 3
	}
	return diff
}

func (h *Handler) getSort(r *http.Request) string {

	return ih.StringFromPath(r, "getStringFromPath", "time")
}

func (h *Handler) getName(r *http.Request) (username string, err error) {

	vars := mux.Vars(r)

	if username = vars["name"]; username == "" {
		return "", re.ErrorInvalidName()
	}

	return
}

func (h *Handler) getNameAndPage(r *http.Request) (page int, username string, err error) {
	vars := mux.Vars(r)

	if username = vars["name"]; username == "" {
		return 0, "", re.ErrorInvalidName()
	}

	if vars["page"] == "" {
		page = 1
	} else {
		if page, err = strconv.Atoi(vars["page"]); err != nil {
			return 0, username, re.ErrorInvalidPage()
		}
		if page < 1 {
			page = 1
		}
	}
	return
}

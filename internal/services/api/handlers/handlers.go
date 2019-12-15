package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	mi "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/middleware"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/router"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
)

type Handlers struct {
	c       *config.Configuration
	user    *UserHandler
	users   *UsersHandler
	game    *GameHandler
	session *SessionHandler
	image   *ImageHandler
}

type HandlerI interface {
	Init(c *config.Configuration, db *database.Input) error
}

// InitWithPostgreSQL apply postgreSQL as database
func (h *Handlers) InitWithPostgreSQL(c *config.Configuration) error {
	return h.Init(c, new(database.Input).InitAsPSQL())
}

// Init open connection to database and put it to all handlers
func (h *Handlers) Init(c *config.Configuration, input *database.Input) error {
	fmt.Println("string:", c.DataBase.ConnectionString)
	h.c = c

	err := input.Database.Open(c.DataBase)
	if err != nil {
		return err
	}

	mi.Init()

	h.user = new(UserHandler)
	h.session = new(SessionHandler)
	h.game = new(GameHandler)
	h.users = new(UsersHandler)
	h.image = new(ImageHandler)
	return InitHandlers(c, input, h.user, h.session,
		h.game, h.users, h.image)

}

// Close connections to darabase of all handlers
func (h *Handlers) Close() error {
	return re.Close(h.user, h.users, h.session, h.game, h.image)
}

// Router return router of api operations
func (h *Handlers) Router() *mux.Router {

	r := mux.NewRouter()

	r.PathPrefix("/swagger").Handler(httpSwagger.Handler(
		httpSwagger.URL("swagger/doc.json"), //The url pointing to API definition"
	))

	r.PathPrefix("/metrics").Handler(promhttp.Handler())

	var api = r.PathPrefix("/api").Subrouter()
	var apiWithAuth = r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/user", h.user.Handle).Methods("OPTIONS", "POST", "DELETE")
	apiWithAuth.HandleFunc("/user", h.user.Handle).Methods("PUT", "GET")

	api.HandleFunc("/session", h.session.Handle).Methods("POST", "OPTIONS")
	apiWithAuth.HandleFunc("/session", h.session.Handle).Methods("DELETE")

	// delete "/avatar/{name}" path
	api.HandleFunc("/avatar/{name}", h.image.Handle).Methods("GET")

	api.HandleFunc("/avatar", h.image.Handle).Methods("OPTIONS")
	apiWithAuth.HandleFunc("/avatar", h.image.Handle).Methods("POST")

	api.HandleFunc("/game", h.game.Handle).Methods("OPTIONS")
	apiWithAuth.HandleFunc("/game", h.game.Handle).Methods("POST")

	api.HandleFunc("/users/{id}", h.users.HandleGetProfile).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/pages/page", h.users.HandleUsersPages).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/pages/amount", h.users.HandleUsersPageAmount).Methods("GET")

	r.PathPrefix("/health").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		rw.Write([]byte("all ok " + server.GetIP()))
	})

	r.PathPrefix("/hard").HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		time.Sleep(7 * time.Second)
		rw.Write([]byte("hard done " + server.GetIP()))
	})

	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	router.Use(r, mi.CORS(h.c.Cors))
	apiWithAuth.Use(mi.Auth(h.c.Cookie, h.c.Auth, h.c.AuthClient))
	return r
}

func InitHandlers(c *config.Configuration, db *database.Input, handlers ...HandlerI) error {
	var err error
	for _, handler := range handlers {
		err = handler.Init(c, db)
		if err != nil {
			break
		}
	}
	return re.Wrap(err)
}

// HistoryRouter return router for history service
/*
func HistoryRouter(handler *api.Handler, cors config.CORS) *mux.Router {
	r := mux.NewRouter()

	var history = r.PathPrefix("/history").Subrouter()

	history.Use(mi.Recover, mi.CORS(cors), mux.CORSMethodMiddleware(r))

	history.HandleFunc("/ws", handler.GameOnline)
	history.Handle("/metrics", promhttp.Handler())
	return r
}
*/

// 128 -> 88

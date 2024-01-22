package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/HeadGardener/effective_mobile/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	minute = time.Minute
)

const (
	personIDParam    = "person_id"
	personIDQuery    = "person_id"
	createdAtQuery   = "created_at"
	limitQuery       = "limit"
	nameQuery        = "name"
	surnameQuery     = "surname"
	ageQuery         = "age"
	genderQuery      = "gender"
	nationalityQuery = "nationality"
)

type PersonService interface {
	Create(ctx context.Context, person *models.Person) (string, error)
	Get(ctx context.Context, filters map[string]any, id, createdAt string, limit int) ([]models.Person, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, fields map[string]any) error
}

type Handler struct {
	log *slog.Logger

	personService PersonService
}

func NewHandler(personService PersonService) *Handler {
	return &Handler{
		log:           slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		personService: personService,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(minute))

	r.Route("/api", func(r chi.Router) {
		r.Post("/", h.createPerson)
		r.Get("/", h.getPersons)
		r.Put("/{person_id}", h.updatePerson)
		r.Delete("/{person_id}", h.deletePerson)
	})

	return r
}

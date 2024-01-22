package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/HeadGardener/effective_mobile/internal/models"
	"github.com/go-chi/chi/v5"
)

var (
	letterRegexp = regexp.MustCompile(`[A-z]$`)
)

type createPersonReq struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type updatePersonRequest struct {
	Name        *string `json:"name"`
	Surname     *string `json:"surname"`
	Patronymic  *string `json:"patronymic"`
	Age         *int8   `json:"age"`
	Gender      *string `json:"gender"`
	Nationality *string `json:"nationality"`
}

func (h *Handler) createPerson(w http.ResponseWriter, r *http.Request) {
	var req createPersonReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.newErrResponse(w, http.StatusBadRequest, "failed while decoding create person req", err)
		return
	}

	if err := req.validate(); err != nil {
		h.newErrResponse(w, http.StatusBadRequest, "failed while validating create person req", err)
		return
	}

	person := &models.Person{
		Name:       req.Name,
		Surname:    req.Surname,
		Patronymic: req.Patronymic,
	}

	id, err := h.personService.Create(r.Context(), person)
	if err != nil {
		h.newErrResponse(w, http.StatusInternalServerError, "failed while creating person", err)
		return
	}

	h.newResponse(w, http.StatusCreated, map[string]any{
		"id": id,
	})
}

func (h *Handler) getPersons(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get(personIDQuery)
	createdAt := r.URL.Query().Get(createdAtQuery)

	if err := validatePersonIDAndCreatedAtQuery(id, createdAt); err != nil {
		h.newErrResponse(w, http.StatusBadRequest, "failed while validation query", err)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get(limitQuery))
	if err != nil {
		h.newErrResponse(w, http.StatusBadRequest, "invalid limit value", err)
		return
	}

	filters, err := queryToMap(r.URL.Query())
	if err != nil {
		h.newErrResponse(w, http.StatusBadRequest, "invalid query params", err)
		return
	}

	persons, err := h.personService.Get(r.Context(), filters, id, createdAt, limit)
	if err != nil {
		h.newErrResponse(w, http.StatusInternalServerError, "failed while getting persons", err)
		return
	}

	h.newResponse(w, http.StatusOK, persons)
}

func (h *Handler) updatePerson(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, personIDParam)

	var req updatePersonRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.newErrResponse(w, http.StatusBadRequest, "failed while decoding update person req", err)
		return
	}

	if err := req.validate(); err != nil {
		h.newErrResponse(w, http.StatusBadRequest, "failed while validating update person req", err)
		return
	}

	fields := req.toMap()

	if err := h.personService.Update(r.Context(), id, fields); err != nil {
		h.newErrResponse(w, http.StatusInternalServerError, "failed while updating person", err)
		return
	}

	h.newResponse(w, http.StatusOK, map[string]any{
		"status": "updated",
	})
}

func (h *Handler) deletePerson(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, personIDParam)

	if err := h.personService.Delete(r.Context(), id); err != nil {
		h.newErrResponse(w, http.StatusInternalServerError, "failed while deleting person", err)
		return
	}

	h.newResponse(w, http.StatusOK, map[string]any{
		"status": "deleted",
	})
}

func (req *createPersonReq) validate() error {
	if (req.Name == "") || !letterRegexp.MatchString(req.Name) {
		return errors.New("invalid name, it must contain only letters and can't be empty")
	}

	if (req.Surname == "") || !letterRegexp.MatchString(req.Surname) {
		return errors.New("invalid surname, it must contain only letters and can't be empty")
	}

	if (req.Patronymic != "") && !letterRegexp.MatchString(req.Surname) {
		return errors.New("invalid patronymic, it must contain only letters and can't be empty")
	}

	return nil
}

func validatePersonIDAndCreatedAtQuery(id, createdAt string) error {
	if id != "" {
		if _, err := uuid.Parse(id); err != nil {
			return fmt.Errorf("invalid id query: %w", err)
		}
	}

	if createdAt != "" {
		if _, err := time.Parse("2006-01-02 15:04:05.000000", createdAt); err != nil {
			return fmt.Errorf("invalid created_at query: %w", err)
		}
	}

	return nil
}

func queryToMap(vals url.Values) (map[string]any, error) {
	var m = make(map[string]any)

	if vals.Has(nameQuery) {
		m[nameQuery] = vals.Get(nameQuery)
	}

	if vals.Has(surnameQuery) {
		m[surnameQuery] = vals.Get(surnameQuery)
	}

	if vals.Has(ageQuery) {
		age, err := strconv.Atoi(vals.Get(ageQuery))
		if err != nil {
			return nil, err
		}

		m[ageQuery] = age
	}

	if vals.Has(genderQuery) {
		m[genderQuery] = vals.Get(genderQuery)
	}

	if vals.Has(nationalityQuery) {
		m[nationalityQuery] = vals.Get(nationalityQuery)
	}

	return m, nil
}

func (req *updatePersonRequest) validate() error {
	if (req.Name != nil) && (!letterRegexp.MatchString(*req.Name) || (*req.Name == "")) {
		return errors.New("invalid name to update, it must contain only letters and can't be empty")
	}

	if (req.Surname != nil) && (!letterRegexp.MatchString(*req.Surname) || (*req.Surname == "")) {
		return errors.New("invalid surname to update, it must contain only letters and can't be empty")
	}

	if (req.Patronymic != nil) && (!letterRegexp.MatchString(*req.Patronymic) || (*req.Patronymic == "")) {
		return errors.New("invalid patronymic to update, it must contain only letters and can't be empty")
	}

	if (req.Age != nil) && (*req.Age < 0 || *req.Age > 120) {
		return errors.New("invalid age, it must be greater than 0 and less than 120")
	}

	if (req.Gender != nil) && (!letterRegexp.MatchString(*req.Gender) || (*req.Gender == "")) {
		return errors.New("invalid gender to update, it must contain only letters and it can't be empty")
	}

	if (req.Nationality != nil) && (!letterRegexp.MatchString(*req.Nationality) || (*req.Nationality == "")) {
		return errors.New("invalid nationality to update, it must contain only letters and can't be empty")
	}

	return nil
}

func (req *updatePersonRequest) toMap() map[string]any {
	var m = make(map[string]any)

	if req.Name != nil {
		m["name"] = *req.Name
	}

	if req.Surname != nil {
		m["surname"] = *req.Surname
	}

	if req.Patronymic != nil {
		m["patronymic"] = *req.Patronymic
	}

	if req.Age != nil {
		m["age"] = *req.Age
	}

	if req.Gender != nil {
		m["gender"] = *req.Gender
	}

	if req.Nationality != nil {
		m["nationality"] = *req.Nationality
	}

	return m
}

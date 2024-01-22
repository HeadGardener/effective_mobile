package storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/HeadGardener/effective_mobile/internal/models"
	"github.com/jmoiron/sqlx"
)

type PersonStorage struct {
	db *sqlx.DB

	debugLogger *slog.Logger
}

func NewPersonStorage(db *sqlx.DB) *PersonStorage {
	return &PersonStorage{
		db:          db,
		debugLogger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	}
}

func (s *PersonStorage) Save(ctx context.Context, person *models.Person) (string, error) {
	start := time.Now()
	if _, err := s.db.ExecContext(ctx, `INSERT INTO persons
    										(id, name, surname, patronymic, age, gender, nationality, created_at)
											VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		person.ID,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
		person.CreatedAt); err != nil {
		return "", err
	}

	s.debugLogger.Debug("saved person", "time", time.Since(start).String(), "person_id", person.ID)
	return person.ID, nil
}

func (s *PersonStorage) GetByID(ctx context.Context, id string) (*models.Person, error) {
	start := time.Now()
	var person models.Person

	if err := s.db.GetContext(ctx, &person, `SELECT * FROM persons WHERE id=$1`, id); err != nil {
		return nil, err
	}

	s.debugLogger.Debug("select person by id", "time", time.Since(start).String(), "person_id", id)

	return &person, nil
}

func (s *PersonStorage) Get(ctx context.Context,
	filters map[string]any, id, createdAt string, limit int) ([]models.Person, error) {
	start := time.Now()
	argID := 1
	args := make([]any, 0)

	var pagPart string
	if id != "" && createdAt != "" {
		pagPart = fmt.Sprintf(`(id, created_at) < ($%d, $%d)`, argID, argID+1)
		args = append(args, id, createdAt)
		argID += 2
	}

	getValues := make([]string, 0)
	for column, value := range filters {
		getValues = append(getValues, fmt.Sprintf("%s=$%d", column, argID))
		args = append(args, value)
		argID++
	}

	var query = strings.Builder{}
	query.WriteString("SELECT * FROM persons ")

	if pagPart != "" || len(getValues) != 0 {
		query.WriteString("WHERE ")
	}

	if pagPart != "" {
		query.WriteString(pagPart)
	}

	if pagPart != "" && len(getValues) != 0 {
		query.WriteString(" AND ")
	}

	if len(getValues) != 0 {
		query.WriteString(strings.Join(getValues, ", "))
	}

	query.WriteString(fmt.Sprintf(" ORDER BY created_at DESC, id DESC LIMIT $%d", argID))
	args = append(args, limit)

	s.debugLogger.Debug("build up query", "query", query.String())

	var persons []models.Person

	if err := s.db.SelectContext(ctx, &persons, query.String(), args...); err != nil {
		return nil, err
	}

	s.debugLogger.Debug("select persons with given filters", "time", time.Since(start).String(),
		"filters", filters)

	return persons, nil
}

func (s *PersonStorage) Update(ctx context.Context, id string, fields map[string]any) error {
	start := time.Now()
	setValues := make([]string, 0)
	args := make([]any, 0)
	argID := 1
	for column, value := range fields {
		setValues = append(setValues, fmt.Sprintf("%s=$%d", column, argID))
		args = append(args, value)
		argID++
	}

	query := fmt.Sprintf(`UPDATE persons SET %s WHERE id=$%d`,
		strings.Join(setValues, ", "), argID)
	args = append(args, id)

	s.debugLogger.Debug("build up query", "query", query)

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	s.debugLogger.Debug("update person by id with given fields", "time", time.Since(start).String(),
		"person_id", id, "fields", fields)

	return nil
}

func (s *PersonStorage) Delete(ctx context.Context, id string) error {
	start := time.Now()
	if _, err := s.db.ExecContext(ctx, `DELETE FROM persons WHERE id=$1`, id); err != nil {
		return err
	}

	s.debugLogger.Debug("delete person", "time", time.Since(start).String(), "person_id", id)

	return nil
}

package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/HeadGardener/effective_mobile/internal/models"
	"github.com/jmoiron/sqlx"
)

type PersonStorage struct {
	db *sqlx.DB
}

func NewPersonStorage(db *sqlx.DB) *PersonStorage {
	return &PersonStorage{db: db}
}

func (s *PersonStorage) Save(ctx context.Context, person *models.Person) (string, error) {
	if _, err := s.db.ExecContext(ctx, `INSERT INTO persons
    										(id, name, surname, patronymic, age, gender, nationality, created_at)
											VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`); err != nil {
		return "", err
	}

	return person.ID, nil
}

func (s *PersonStorage) GetByID(ctx context.Context, id string) (*models.Person, error) {
	var person models.Person

	if err := s.db.GetContext(ctx, &person, `SELECT * FROM persons WHERE id=$1`, id); err != nil {
		return nil, err
	}

	return &person, nil
}

func (s *PersonStorage) Get(ctx context.Context,
	filters map[string]any, id, createdAt string, limit int) ([]models.Person, error) {
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

	if pagPart != "" && len(getValues) != 0 {
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

	var persons []models.Person

	if err := s.db.SelectContext(ctx, &persons, query.String(), args...); err != nil {
		return nil, err
	}

	return persons, nil
}

func (s *PersonStorage) Update(ctx context.Context, id string, fields map[string]any) error {
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

	if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (s *PersonStorage) Delete(ctx context.Context, id string) error {
	if _, err := s.db.ExecContext(ctx, `DELETE FROM persons WHERE id=$1`, id); err != nil {
		return err
	}

	return nil
}

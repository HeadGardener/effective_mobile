package services

import (
	"context"
	"time"

	"github.com/HeadGardener/effective_mobile/internal/models"
	"github.com/google/uuid"
)

type PersonStorage interface {
	Save(ctx context.Context, person *models.Person) (string, error)
	GetByID(ctx context.Context, id string) (*models.Person, error)
	Get(ctx context.Context, filters map[string]any, id, createdAt string, limit int) ([]models.Person, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, fields map[string]any) error
}

type PersonDataProvider interface {
	GetAge(ctx context.Context, name string) (int8, error)
	GetGender(ctx context.Context, name string) (string, error)
	GetNationality(ctx context.Context, name string) (string, error)
}

type PersonService struct {
	personStorage      PersonStorage
	personDataProvider PersonDataProvider
}

func NewPersonService(personStorage PersonStorage, personDataProvider PersonDataProvider) *PersonService {
	return &PersonService{
		personStorage:      personStorage,
		personDataProvider: personDataProvider,
	}
}

func (s *PersonService) Create(ctx context.Context, person *models.Person) (string, error) {
	age, err := s.personDataProvider.GetAge(ctx, person.Name)
	if err != nil {
		return "", err
	}
	person.Age = age

	gender, err := s.personDataProvider.GetGender(ctx, person.Name)
	if err != nil {
		return "", err
	}
	person.Gender = gender

	nationality, err := s.personDataProvider.GetNationality(ctx, person.Name)
	if err != nil {
		return "", err
	}
	person.Nationality = nationality

	person.ID = uuid.NewString()
	person.CreatedAt = time.Now()

	return s.personStorage.Save(ctx, person)
}

func (s *PersonService) Get(ctx context.Context, filters map[string]any, id, createdAt string, limit int) ([]models.Person, error) {
	return s.personStorage.Get(ctx, filters, id, createdAt, limit)
}

func (s *PersonService) Delete(ctx context.Context, id string) error {
	if _, err := s.personStorage.GetByID(ctx, id); err != nil {
		return err
	}

	return s.personStorage.Delete(ctx, id)
}

func (s *PersonService) Update(ctx context.Context, id string, fields map[string]any) error {
	if _, err := s.personStorage.GetByID(ctx, id); err != nil {
		return err
	}

	return s.personStorage.Update(ctx, id, fields)
}

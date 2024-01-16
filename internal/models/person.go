package models

import "time"

type Person struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Surname     string    `db:"surname"`
	Patronymic  string    `db:"patronymic"`
	Age         int8      `db:"age"`
	Gender      string    `db:"gender"`
	Nationality string    `db:"nationality"`
	CreatedAt   time.Time `db:"created_at"`
}

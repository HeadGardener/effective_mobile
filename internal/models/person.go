package models

type Person struct {
	Name        string `db:"name"`
	Surname     string `db:"surname"`
	Patronymic  string `db:"patronymic"`
	Age         int8   `db:"age"`
	Gender      string `db:"gender"`
	Nationality string `db:"nationality"`
}

-- +goose Up
-- +goose StatementBegin
CREATE TABLE persons
(
    id          uuid PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    surname     VARCHAR(255) NOT NULL,
    patronymic  VARCHAR(255) NOT NULL,
    age         int          NOT NULL,
    gender      VARCHAR(50)  NOT NULL,
    nationality VARCHAR(50)  NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE persons;
-- +goose StatementEnd

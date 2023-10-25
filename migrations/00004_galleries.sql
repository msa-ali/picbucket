-- +goose Up
-- +goose StatementBegin
CREATE TABLE galleries (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    title TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE galleries;
-- +goose StatementEnd

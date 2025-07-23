-- +goose Up
-- +goose StatementBegin
CREATE TABLE syllabus (
    id BIGSERIAL PRIMARY KEY,
    week TEXT,
    topic TEXT
);

CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    task TEXT,
    description TEXT,
    solution TEXT,
    hint1 TEXT,
    hint2 TEXT,
    hint3 TEXT,
    difficulty INT,
    solved INT
);

CREATE TABLE statistics(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    easy INT,
    medium INT,
    hard INT,
    total INT
);

CREATE TABLE notifications(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    enabled BOOL,
    time_24 TIMESTAMP,
    days INT[]
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE syllabus;
DROP TABLE tasks;
DROP TABLE statistics;
DROP TABLE notifications;
-- +goose StatementEnd

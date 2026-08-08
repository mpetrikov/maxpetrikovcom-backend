CREATE TABLE roles (
    id SMALLINT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO roles (id, name)
VALUES
    (1, 'student'),
    (2, 'admin');
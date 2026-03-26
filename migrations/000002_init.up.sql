CREATE SCHEMA todolist;

CREATE TABLE todolist.users (
    id               SERIAL                PRIMARY KEY,
    version          BIGINT       NOT NULL DEFAULT 1,
    full_name        VARCHAR(100) NOT NULL CHECK (char_length(full_name) BETWEEN 3 AND 100),
    phone_number     VARCHAR(15)           CHECK (
        phone_number ~ '^\+?[0-9]+$'
        AND
        char_length(phone_number) BETWEEN 10 AND 15
    )

);

CREATE TABLE todolist.tasks (
    id               SERIAL                PRIMARY KEY,
    version          BIGINT       NOT NULL DEFAULT 1,
    title            VARCHAR(255) NOT NULL CHECK (char_length(title) BETWEEN 1 AND 255),
    description      VARCHAR(1000),
    is_completed     BOOLEAN       NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ   NOT NULL,
    completed_at     TIMESTAMPTZ,  

    CHECK (
        is_completed=FALSE AND completed_at IS NULL
        OR
        is_completed=TRUE AND completed_at IS NOT NULL AND completed_at >= created_at
    ),

    author_user_id   INTEGER       NOT NULL REFERENCES todolist.users(id)
);
CREATE TABLE IF NOT EXISTS accounts (
    id bytea NOT NULL,
    email varchar NOT NULL,
    password varchar NOT NULL,
    verification_code char(10) NOT NULL,
    verified_at timestamp,
    created_at timestamp NOT NULL,
    updated_at timestamp,

    PRIMARY KEY(id)
)
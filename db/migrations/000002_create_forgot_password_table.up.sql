CREATE TABLE IF NOT EXISTS forgot_passwords (
    id bytea NOT NULL,
    user_id bytea NOT NULL,
    verification_code char(32) NOT NULL,
    valid boolean NOT NULL DEFAULT 1,
    created_at timestamp NOT NULL,
    updated_at timestamp,
    expired_at timestamp NOT NULL,

    PRIMARY KEY(id)
)
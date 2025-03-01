CREATE TABLE IF NOT EXISTS users (
    guid          uuid PRIMARY KEY NOT NULL,
    email         varchar(20) NOT NULL,
    password_hash varchar(60) NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT NOW(),
    updated_at    timestamptz NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX users_email_uniq_idx ON users (email);

COMMENT ON TABLE users IS 'Пользователи';
COMMENT ON COLUMN users.guid IS 'GUID';
COMMENT ON COLUMN users.email IS 'Email';
COMMENT ON COLUMN users.created_at IS 'Дата создания пользователя';
COMMENT ON COLUMN users.updated_at IS 'Дата изменения пользователя';
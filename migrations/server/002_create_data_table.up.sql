CREATE TABLE IF NOT EXISTS items (
    guid           uuid PRIMARY KEY NOT NULL,
    user_guid      uuid NOT NULL REFERENCES users (guid) ON DELETE RESTRICT,
    encrypted_data text NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT NOW(),
    updated_at     timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX items_user_guid_idx ON items (user_guid);

COMMENT ON TABLE items IS 'Данные (пароли, логины, тексты и пр.)';
COMMENT ON COLUMN items.guid IS 'GUID';
COMMENT ON COLUMN items.user_guid IS 'Владелец';
COMMENT ON COLUMN items.encrypted_data IS 'Сами данные в зашифрованном виде';
COMMENT ON COLUMN items.created_at IS 'Дата создания записи';
COMMENT ON COLUMN items.updated_at IS 'Дата изменения записи';
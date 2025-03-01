CREATE TABLE IF NOT EXISTS data (
    guid           uuid PRIMARY KEY NOT NULL,
    user_guid      uuid NOT NULL REFERENCES users (guid) ON DELETE RESTRICT,
    encrypted_data text NOT NULL,
    created_at     timestamptz NOT NULL DEFAULT NOW(),
    updated_at     timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX data_user_guid_idx ON data (user_guid);

COMMENT ON TABLE data IS 'Данные (пароли, логины, тексты и пр.)';
COMMENT ON COLUMN data.guid IS 'GUID';
COMMENT ON COLUMN data.user_guid IS 'Владелец';
COMMENT ON COLUMN data.encrypted_data IS 'Сами данные в зашифрованном виде';
COMMENT ON COLUMN data.created_at IS 'Дата создания записи';
COMMENT ON COLUMN data.updated_at IS 'Дата изменения записи';
-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at_column() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now() at time zone 'utc';
RETURN NEW;
END;
$$ language 'plpgsql';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE campaigns
(
    id          uuid                              DEFAULT gen_random_uuid() NOT NULL,
    name        TEXT                     NOT NULL,
    description TEXT                     NOT NULL,
    adserver_id uuid NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc'),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now() at time zone 'utc')
);

CREATE UNIQUE INDEX campaigns__name__uidx ON campaigns USING btree (name);
CREATE UNIQUE INDEX campaigns__adserver_id__uidx ON campaigns USING btree (adserver_id);

CREATE TRIGGER update_campaign_updated_at
    BEFORE
        UPDATE
    ON campaigns
    FOR EACH ROW
EXECUTE PROCEDURE set_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE campaigns;
-- +goose StatementEnd

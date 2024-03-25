-- +goose Up
-- +goose StatementBegin
create table pickup_points(
    id BIGSERIAL PRIMARY KEY NOT NULL
    name TEXT NOT NULL DEFAULT ''
    address TEXT NOT NULL DEFAULT ''
    phone_number TEXT NOT NULL DEFAULT ''
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table pickup_points
-- +goose StatementEnd

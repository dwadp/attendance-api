-- +goose Up
-- +goose StatementBegin
create table "public"."employees" (
    "id" serial primary key not null,
    "name" varchar(50) not null,
    "phone" varchar(25) not null,
    "created_at" timestamp null default NOW(),
    "updated_at" timestamp null default NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "public"."employees";
-- +goose StatementEnd

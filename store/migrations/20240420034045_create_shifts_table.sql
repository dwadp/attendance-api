-- +goose Up
-- +goose StatementBegin
create table "public"."shifts" (
    "id" serial primary key not null,
    "name" varchar(20) not null,
    "in" time not null,
    "out" time not null,
    "is_default" boolean not null default false,
    "created_at" timestamp null default NOW(),
    "updated_at" timestamp null default NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "public"."shifts";
-- +goose StatementEnd

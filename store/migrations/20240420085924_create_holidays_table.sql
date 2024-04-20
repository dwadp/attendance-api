-- +goose Up
-- +goose StatementBegin
create table "public"."holidays" (
    "id" serial primary key not null,
    "name" varchar(50) not null,
    "type" int not null,
    "weekday" int null,
    "date" date null,
    "created_at" timestamp null default NOW(),
    "updated_at" timestamp null default NOW()
);

CREATE INDEX ON "public"."holidays" (type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "public"."holidays";
-- +goose StatementEnd

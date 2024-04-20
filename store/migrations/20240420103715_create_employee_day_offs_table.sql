-- +goose Up
-- +goose StatementBegin
create table "public"."employee_day_offs" (
    "id" serial primary key not null,
    "employee_id" integer not null,
    "description" text not null,
    "date" date not null,
    "created_at" timestamp null default NOW(),
    CONSTRAINT fk_employee_day_offs_employee_id
        FOREIGN KEY(employee_id) REFERENCES employees(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "public"."employee_day_offs";
-- +goose StatementEnd

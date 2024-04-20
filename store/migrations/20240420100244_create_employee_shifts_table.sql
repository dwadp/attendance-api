-- +goose Up
-- +goose StatementBegin
create table "public"."employee_shifts" (
    "id" serial primary key not null,
    "employee_id" integer not null,
    "shift_id" integer not null,
    "date" date not null,
    "created_at" timestamp not null default NOW(),
    CONSTRAINT fk_employee_shifts_employee_id
        FOREIGN KEY(employee_id) REFERENCES employees(id),
    CONSTRAINT fk_employee_shifts_shift_id
        FOREIGN KEY(shift_id) REFERENCES shifts(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "public"."employee_shifts";
-- +goose StatementEnd

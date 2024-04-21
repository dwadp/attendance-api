-- +goose Up
-- +goose StatementBegin
create table "public"."attendances" (
    "id" serial primary key not null,
    "employee_id" integer not null,
    "shift_id" integer null,
    "shift_name" varchar(20) not null,
    "shift_in" timestamp not null,
    "shift_out" timestamp not null,
    "clock_in" timestamp null,
    "clock_out" timestamp null,
    "clock_in_status" varchar(20) null,
    "clock_out_status" varchar(20) null,
    "date" date not null,
    "created_at" timestamp null default NOW(),
    "updated_at" timestamp null default NOW(),
    CONSTRAINT fk_attendances_employee_id
        FOREIGN KEY(employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_attendances_shift_id
        FOREIGN KEY(shift_id) REFERENCES shifts(id) ON DELETE SET NULL
);

CREATE INDEX ON "public"."attendances" (clock_in_status);
CREATE INDEX ON "public"."attendances" (clock_out_status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "public"."attendances";
-- +goose StatementEnd

-- Write your migrate up statements here

create table users
(
    id       uuid default gen_random_uuid() not null constraint users_pk unique primary key,
    login    varchar                        not null unique,
    password varchar                        not null
);

create table balances
(
    user_id    uuid constraint balances_pk unique constraint balances_users_null_fk references users (id),
    current    integer default 0,
    withdrawal integer default 0,
    created_at timestamptz default current_timestamp not null
);

create table orders
(
    id         uuid        default gen_random_uuid() not null constraint orders_pk unique primary key,
    number     varchar                               not null constraint orders_number_uniqk unique,
    status     varchar,
    accrual    integer,
    user_id    uuid                                  not null constraint orders_users_null_fk references users (id),
    created_at timestamptz default current_timestamp
);

create table if not exists withdrawals
(
    id         uuid        default gen_random_uuid() not null constraint withdrawals_pk primary key,
    "order"    varchar                               not null constraint withdrawals_order_uniqk unique,
    sum        integer                               not null,
    user_id    uuid                                  not null constraint withdrawals_user_id_fk references users (id),
    created_at timestamptz default current_timestamp not null
);

---- create above / drop below ----

drop table withdrawals;

drop table orders;

drop table balances;

drop table users;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

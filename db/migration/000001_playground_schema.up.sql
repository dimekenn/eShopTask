create table users
(
    id       bigserial primary key,
    username varchar not null,
    password text    not null
);

create table roles
(
    id   serial primary key,
    name varchar not null
);

create table user_roles
(
    user_id int not null,
    role_id int not null,
    constraint fk_user_id foreign key (user_id) references users (id),
    constraint fk_role_id foreign key (role_id) references roles (id)
);

create table products
(
    id          serial primary key,
    name        varchar,
    description text,
    count       int
);

create table carts
(
    id      bigserial primary key,
    user_id int not null,
    constraint fk_user_id foreign key (user_id) references users (id)
);

create table cart_products
(
    id         bigserial primary key,
    cart_id    int not null,
    product_id int not null,
    count      int not null,
    constraint fk_cart_id foreign key (cart_id) references carts (id),
    constraint fk_product_id foreign key (product_id) references products (id)
);

insert into roles (id, name) vales (1, 'manager'), (2, 'user');
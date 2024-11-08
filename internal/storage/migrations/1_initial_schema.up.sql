CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    login VARCHAR not null,
    password VARCHAR not null
);
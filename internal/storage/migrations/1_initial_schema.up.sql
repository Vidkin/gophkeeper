CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR not null,
    password VARCHAR not null
);

CREATE TABLE bank_cards (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    owner TEXT NOT NULL,
    card_number TEXT NOT NULL,
    expiration_date TEXT NOT NULL,
    cvv TEXT NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE user_credentials (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    login TEXT NOT NULL,
    password TEXT NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    text TEXT NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE files (
   id SERIAL PRIMARY KEY,
   user_id INT NOT NULL,
   bucket_name VARCHAR(255) NOT NULL,
   file_name VARCHAR(255) NOT NULL,
   file_size BIGINT NOT NULL,
   description VARCHAR(255),
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);

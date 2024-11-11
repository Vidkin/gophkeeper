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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE files (
   id SERIAL PRIMARY KEY,
   user_id INT NOT NULL,
   bucket_name VARCHAR(255) NOT NULL,
   file_name VARCHAR(255) NOT NULL,
   file_type VARCHAR(50) NOT NULL, -- Тип файла (например, "text/plain", "image/jpeg" и т.д.)
   file_size BIGINT NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);

-- +goose Up
-- +goose StatementBegin
-- Создание таблицы пользователей
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       coins INT NOT NULL DEFAULT 1000 CHECK (coins >= 0)
);

CREATE TABLE merch (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) UNIQUE NOT NULL,
                       price INT NOT NULL CHECK (price > 0)
);

CREATE TABLE transactions (
                              id SERIAL PRIMARY KEY,
                              from_user INT REFERENCES users(id),
                              to_user INT REFERENCES users(id) NOT NULL,
                              amount INT NOT NULL CHECK (amount > 0),
                              created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE purchases (
                           id SERIAL PRIMARY KEY,
                           user_id INT REFERENCES users(id) NOT NULL,
                           merch_id INT REFERENCES merch(id) NOT NULL,
                           quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
                           created_at TIMESTAMP DEFAULT NOW()
);

-- Предзаполнение товаров
INSERT INTO merch (name, price) VALUES
                                    ('t-shirt', 80),
                                    ('cup', 20),
                                    ('book', 50),
                                    ('pen', 10),
                                    ('powerbank', 200),
                                    ('hoody', 300),
                                    ('umbrella', 200),
                                    ('socks', 10),
                                    ('wallet', 50),
                                    ('pink-hoody', 500);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE purchases;
DROP TABLE transactions;
DROP TABLE merch;
DROP TABLE users;
-- +goose StatementEnd

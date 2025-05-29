-- +goose Up
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    userid TEXT NOT NULL,
    price FLOAT NOT NULL,
    quantity INT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE orders;

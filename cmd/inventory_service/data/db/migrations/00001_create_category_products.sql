-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS category(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    tags TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    price_in_cents INTEGER NOT NULL CHECK (price_in_cents >= 0),
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS category_products_relation(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,

    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES category(id),
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS category_products_relation;

-- +goose StatementEnd
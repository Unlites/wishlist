CREATE SCHEMA IF NOT EXISTS wishlist;

CREATE TABLE IF NOT EXISTS wishlist.users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    info TEXT
);

CREATE TABLE IF NOT EXISTS wishlist.wishes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES wishlist.users (id),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    is_reserved BOOLEAN DEFAULT FALSE NOT NULL,
    reserved_by INTEGER REFERENCES wishlist.users (id),
    price INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

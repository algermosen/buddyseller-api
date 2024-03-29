CREATE TABLE
    IF NOT EXISTS products (
        id serial PRIMARY KEY,
        name character(30) NOT NULL,
        description text NOT NULL,
        sku text NOT NULL UNIQUE,
        price text NOT NULL,
        stock integer NOT NULL
    );

CREATE TABLE
    IF NOT EXISTS users (
        id serial PRIMARY KEY,
        name text NOT NULL,
        code text NOT NULL UNIQUE,
        email text NOT NULL UNIQUE,
        password text NOT NULL
    );
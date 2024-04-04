CREATE TYPE order_status AS ENUM (
  'pending',
  'shipped',
  'delivered',
  'cancelled'
);

CREATE TABLE
    IF NOT EXISTS products (
        id serial PRIMARY KEY,
        name text NOT NULL,
        description text NOT NULL,
        sku text NOT NULL UNIQUE,
        price numeric NOT NULL,
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

CREATE TABLE
    IF NOT EXISTS orders (
        id serial PRIMARY KEY,
        status order_status NOT NULL,
        total_amount numeric NOT NULL,
        tax numeric NOT NULL,
        created timestamp DEFAULT NOW(),
        shipped timestamp,
        cancelled timestamp,
        delivered timestamp,
        client_name text,
        client_email text,
        client_phone text,
        note text,
        cancellation_reason text
    );

CREATE TABLE
    IF NOT EXISTS order_items (
        id serial PRIMARY KEY,
        unit_price numeric NOT NULL,
        quantity integer NOT NULL,
        order_id integer NOT NULL REFERENCES orders(id),
        product_id integer NOT NULL REFERENCES products(id)
    );

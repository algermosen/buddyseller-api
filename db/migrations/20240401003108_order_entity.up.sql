CREATE TYPE IF NOT EXISTS order_status AS ENUM ('pending', 'shipped', 'delivered', 'cancelled');

CREATE TABLE
    IF NOT EXISTS orders (
        id serial PRIMARY KEY,
        status order_status NOT NULL,
        total_amount numeric NOT NULL,
        tax numeric NOT NULL,
        created timestamp DEFAULT NOW (),
        shipped timestamp,
        cancelled timestamp,
        delivered timestamp,
        client_name text,
        client_email text,
        client_phone text,
        note text,
        cancellation_reason text
    );
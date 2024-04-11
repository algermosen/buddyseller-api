CREATE TYPE order_status AS ENUM ('pending', 'shipped', 'delivered', 'cancelled');

CREATE TABLE
    IF NOT EXISTS orders (
        id serial PRIMARY KEY,
        status order_status NOT NULL DEFAULT 'pending',
        total_amount numeric NOT NULL,
        tax numeric NOT NULL,
        user_id integer NOT NULL REFERENCES users(id),
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
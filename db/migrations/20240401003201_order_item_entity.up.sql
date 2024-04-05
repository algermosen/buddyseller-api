CREATE TABLE
    IF NOT EXISTS order_items (
        id serial PRIMARY KEY,
        unit_price numeric NOT NULL,
        quantity integer NOT NULL,
        order_id integer NOT NULL REFERENCES orders(id),
        product_id integer NOT NULL REFERENCES products(id)
    );
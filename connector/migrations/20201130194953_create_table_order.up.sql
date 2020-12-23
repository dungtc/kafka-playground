-- Create table order if not exists
CREATE TABLE IF NOT EXISTS orders(
  id SERIAL PRIMARY KEY,
  item_id integer NOT NULL,
  quantity integer NOT NULL,
  created_at timestamp NOT NULL DEFAULT current_timestamp,
  updated_at timestamp NOT NULL DEFAULT current_timestamp,
  deleted_at timestamp
);
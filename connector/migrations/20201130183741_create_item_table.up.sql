-- Create database
-- CREATE DATABASE sale;

-- Create table if not exists
CREATE TABLE IF NOT EXISTS items(
  id              SERIAL PRIMARY KEY,
  total_stock_value integer NOT NULL,
  current_stock_value integer NOT NULL,
  selling_price decimal NOT NULL DEFAULT '0',
  created_at timestamp NOT NULL DEFAULT current_timestamp,
  updated_at timestamp NOT NULL DEFAULT current_timestamp,
  deleted_at timestamp
);
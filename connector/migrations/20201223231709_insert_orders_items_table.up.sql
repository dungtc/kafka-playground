-- INSERT into items table
INSERT INTO items(id,total_stock_value, current_stock_value, selling_price) values(1,100,100, 10000);
INSERT INTO items(id,total_stock_value, current_stock_value, selling_price) values(2,200,200, 10000);


-- INSERT into order table
INSERT INTO orders(item_id, quantity) values(1, 10);
INSERT INTO orders(item_id, quantity) values(1, 10);
INSERT INTO orders(item_id, quantity) values(1, 10);

UPDATE items SET current_stock_value = 70 WHERE id = 1;

INSERT INTO orders(item_id, quantity) values(2, 20);
INSERT INTO orders(item_id, quantity) values(2, 20);
INSERT INTO orders(item_id, quantity) values(2, 20);

UPDATE items SET current_stock_value = 170 WHERE id = 2;
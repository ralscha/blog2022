CREATE TABLE app_user (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(200) NOT NULL UNIQUE
);


INSERT INTO app_user (username, email) VALUES
  ('demo1', 'demo1@example.com'),
  ('demo2', 'demo2@example.com'),
  ('demo3', 'demo3@example.com'),
  ('demo4', 'demo4@example.com'),
  ('demo5', 'demo5@example.com'),
  ('demo6', 'demo6@example.com'),
  ('demo7', 'demo7@example.com'),
  ('demo8', 'demo8@example.com'),
  ('demo9', 'demo9@example.com'),
  ('demo10', 'demo10@example.com');
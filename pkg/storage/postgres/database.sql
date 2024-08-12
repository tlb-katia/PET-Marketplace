CREATE TABLE users (
    id bigserial not null primary key,
    login VARCHAR(255) not null unique,
    password VARCHAR(255) not null
);

CREATE TABLE adverts (
    id  bigserial not null primary key,
    user_id INT,
    header VARCHAR(64),
    text VARCHAR(512),
    address VARCHAR(64),
    image_url VARCHAR(64),
    price REAL,
    datetime TIMESTAMP
);


ALTER TABLE adverts ALTER COLUMN user_id TYPE BIGINT;

-- Добавьте или измените внешний ключ
ALTER TABLE adverts DROP CONSTRAINT IF EXISTS fk_user;
ALTER TABLE adverts ADD CONSTRAINT fk_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

drop table users;
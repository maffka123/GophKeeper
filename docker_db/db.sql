CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login varchar(40) UNIQUE,
    password varchar(256),
    created_at timestamp DEFAULT current_timestamp
);


CREATE TABLE IF NOT EXISTS secrets  (
    id UUID PRIMARY KEY,
    user_id bigint,
    data jsonb,
    metadata varchar(100),
    change_date timestamp DEFAULT current_timestamp,
    synchronized boolean,
    FOREIGN KEY(user_id) REFERENCES users(id)
);

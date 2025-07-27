CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    feeling VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP NOT NULL
);


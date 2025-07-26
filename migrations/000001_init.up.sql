CREATE TABLE IF NOT EXISTS subscriptions (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    service_name VARCHAR(64) NOT NULL,
    price int NOT NULL,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP
);


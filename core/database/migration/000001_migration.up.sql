CREATE TABLE "order" (
    id BIGSERIAL PRIMARY KEY,
    description VARCHAR(50) NOT NULL,
    transaction_date TIMESTAMP NOT NULL,
    transaction_value FLOAT NOT NULL
);

CREATE TABLE subscriptions (
    user_id INTEGER NOT NULL,
    pair VARCHAR(256) NOT NULL,
    PRIMARY KEY (user_id, pair)
)
CREATE TABLE IF NOT EXISTS
checkout_histories (
    id VARCHAR(16) PRIMARY KEY,
    user_id VARCHAR(16) NOT NULL,
    product_details JSONB NOT NULL,
    paid INT NOT NULL,
    change INT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE checkout_histories
	ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS checkout_histories_user_id
	ON checkout_histories USING HASH(user_id);
CREATE INDEX IF NOT EXISTS checkout_histories_created_at_desc
	ON checkout_histories(created_at DESC);
CREATE INDEX IF NOT EXISTS checkout_histories_created_at_asc
	ON checkout_histories(created_at ASC);
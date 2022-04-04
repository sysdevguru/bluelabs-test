CREATE TABLE wallets (
    id SERIAL,
    user_id bigint NOT NULL,
    balance float
);

ALTER TABLE wallets
    ADD CONSTRAINT wallet_id_pkey PRIMARY KEY (id);

ALTER TABLE wallets
    ADD CONSTRAINT user_id_unique UNIQUE (user_id);

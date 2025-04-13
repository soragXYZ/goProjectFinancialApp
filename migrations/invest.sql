DROP TABLE IF EXISTS invest;
CREATE TABLE invest (
    invest_id INT NOT NULL,
    account_id INT NOT NULL,
    type_id INT NOT NULL,
    invest_label VARCHAR(255) NOT NULL,
    invest_code VARCHAR(255) NOT NULL,
    invest_code_type VARCHAR(255) NOT NULL,
    stock_symbol VARCHAR(255) NOT NULL,
    quantity FLOAT NOT NULL,
    unit_price FLOAT NOT NULL,
    unit_value FLOAT NOT NULL,
    valuation FLOAT NOT NULL,
    diff FLOAT NOT NULL,
    diff_percent FLOAT NOT NULL,
    last_update VARCHAR(255) NOT NULL,

    PRIMARY KEY (`invest_id`),
    FOREIGN KEY (`account_id`) REFERENCES bankAccount(`account_id`)
);

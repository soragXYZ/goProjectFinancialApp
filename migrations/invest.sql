DROP TABLE IF EXISTS invest;
CREATE TABLE invest (
    invest_id INT NOT NULL,
    account_id INT NOT NULL,
    type_id INT NOT NULL,
    invest_label DATETIME NOT NULL,
    invest_code FLOAT NOT NULL,
    invest_code_type VARCHAR(255) NOT NULL,
    stock_symbol VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    unit_price INT NOT NULL,
    unit_value INT NOT NULL,
    valuation INT NOT NULL,
    diff INT NOT NULL,
    diff_percent INT NOT NULL,
    last_update VARCHAR(255) NOT NULL,

    PRIMARY KEY (`invest_id`),
    FOREIGN KEY (`account_id`) REFERENCES bankAccount(`bank_id`)
);

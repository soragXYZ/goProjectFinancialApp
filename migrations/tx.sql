DROP TABLE IF EXISTS tx;
CREATE TABLE tx (
    tx_id INT NOT NULL,
    bank_id INT NOT NULL,
    tx_datetime DATETIME NOT NULL,
    tx_value FLOAT NOT NULL,
    tx_type VARCHAR(255) NOT NULL,
    original_wording VARCHAR(255) NOT NULL,

    PRIMARY KEY (`tx_id`),
    FOREIGN KEY (`bank_id`) REFERENCES bankAccount(`bank_id`)
);

DROP TABLE IF EXISTS historyInvest;
CREATE TABLE historyInvest (
    history_invest_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    invest_id INT NOT NULL,
    valuation FLOAT NOT NULL,
    date_valuation VARCHAR(255) NOT NULL,

    PRIMARY KEY (`history_invest_id`),
    FOREIGN KEY (`invest_id`) REFERENCES invest(`invest_id`)
);

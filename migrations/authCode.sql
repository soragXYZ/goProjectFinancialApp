DROP TABLE IF EXISTS authCode;
CREATE TABLE authCode (
    auth_code VARCHAR(255) NOT NULL,
    code_type VARCHAR(255) NOT NULL,
    access_type VARCHAR(255) NOT NULL,
    expires_in INT NOT NULL,
    PRIMARY KEY (`auth_code`)
);

DROP TABLE IF EXISTS authToken;
CREATE TABLE authToken (
    auth_token VARCHAR(255) NOT NULL,
    token_type VARCHAR(255) NOT NULL,
    id_user INT NOT NULL,
    expires_in INT NOT NULL,
    PRIMARY KEY (`Id_user`)
);

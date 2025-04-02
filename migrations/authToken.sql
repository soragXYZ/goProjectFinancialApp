DROP TABLE IF EXISTS authToken;
CREATE TABLE authToken (
    auth_token VARCHAR(255),
    token_type VARCHAR(255),
    id_user INT,
    expires_in INT,
    PRIMARY KEY (`Id_user`)
);

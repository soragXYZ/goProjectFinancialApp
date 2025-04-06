DROP TABLE IF EXISTS authToken;
CREATE TABLE authToken (
    auth_token VARCHAR(255) NOT NULL,
    id_user INT NOT NULL,
    PRIMARY KEY (`Id_user`)
);

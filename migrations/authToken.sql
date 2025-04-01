-- 

DROP TABLE IF EXISTS authToken;
CREATE TABLE authToken (
    auth_token VARCHAR(255),
    token_type VARCHAR(255),
    id_user INT,
    expires_in INT,
    PRIMARY KEY (`Id_user`)
);


	-- Auth_token string `json:"auth_token"`
	-- Token_type string `json:"type"`
	-- Id_user    int    `json:"id_user"`
	-- Expires_in int    `json:"expires_in"`

--       id         INT AUTO_INCREMENT NOT NULL,
--   title      VARCHAR(128) NOT NULL,
--   artist     VARCHAR(255) NOT NULL,
--   price      DECIMAL(5,2) NOT NULL,
--   PRIMARY KEY (`id`)
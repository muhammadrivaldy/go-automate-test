CREATE TABLE mst_users (
    id BIGINT NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    address TEXT NOT NULL,
    password BLOB NOT NULL,
    PRIMARY KEY (id)
)
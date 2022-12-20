CREATE DATABASE database;

CREATE TABLE users
(
    ID int NOT NULL,
    login   TEXT    NOT NULL,
    password    TEXT    NOT NULL,
    role TEXT   NOT NULL DEFAULT '',
    PRIMARY KEY (ID)
);
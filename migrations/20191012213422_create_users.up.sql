CREATE TABLE users (
    accountId bigserial not null primary key,
    firstName varchar,
    secondName varchar,
    userName varchar not null,
    email varchar not null unique,
    encryptPassword varchar not null,
    avatar bytea
);
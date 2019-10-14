CREATE TABLE companies (
    id bigserial not null primary key,
    name varchar
);

CREATE TABLE managers (
    id bigserial not null primary key,
    accountId bigserial references users,
    registrationDate timestamp,
    location varchar,
    companyId bigserial references  companies
);
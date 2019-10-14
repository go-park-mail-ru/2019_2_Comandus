CREATE TABLE specialities (
    id bigserial not null primary key,
    name varchar
);

CREATE TABLE freelancers (
    id bigserial not null primary key,
    accountId bigserial not null references users,
    registrationDate timestamp,
    country varchar,
    city varchar,
    address varchar,
    phone varchar,
    tagLine varchar,
    overview varchar,
    experienceLevelId bigserial,
    specialityId bigserial --references specialities
);

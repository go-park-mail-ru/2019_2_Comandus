CREATE TABLE jobs (
    id bigserial not null primary key,
    managerId bigserial not null references managers,
    title varchar not null,
    description varchar not null,
    files varchar,
    specialityId bigserial references specialities,
    experienceLevelId bigserial,
    paymentAmount float8,
    country varchar,
    city varchar,
    jobTypeId bigserial
);
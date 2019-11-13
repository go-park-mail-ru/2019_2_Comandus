package create

import "database/sql"

func CreateTables(db *sql.DB) error {
	companiesQuery := `CREATE TABLE IF NOT EXISTS companies (
		id bigserial not null primary key,
		companyName varchar not null,
		site varchar,
		tagLine varchar,
		description varchar,
		country varchar,
		city varchar,
		address varchar,
		phone varchar
	);`
	if _, err := db.Exec(companiesQuery); err != nil {
		return err
	}

	usersQuery := `CREATE TABLE IF NOT EXISTS users (
		accountId bigserial not null primary key,
		firstName varchar,
		secondName varchar,
		userName varchar not null,
		email varchar not null unique,
		encryptPassword varchar not null,
		avatar bytea,
		registrationDate timestamp,
		userType varchar not null
	);`
	if _, err := db.Exec(usersQuery); err != nil {
		return err
	}

	managersQuery := `CREATE TABLE IF NOT EXISTS managers (
		id bigserial not null primary key,
		accountId bigserial references users,
		location varchar,
		companyId bigserial references companies
	);`
	if _, err := db.Exec(managersQuery); err != nil {
		return err
	}

	freelancersQuery := `CREATE TABLE IF NOT EXISTS freelancers (
		id bigserial not null primary key,
		accountId bigserial not null references users,
		country varchar,
		city varchar,
		address varchar,
		phone varchar,
		tagLine varchar,
		overview varchar,
		experienceLevelId bigserial,
		specialityId bigserial --references specialities
	);`
	if _, err := db.Exec(freelancersQuery); err != nil {
		return err
	}

	jobsQuery := `CREATE TABLE IF NOT EXISTS jobs (
		id bigserial not null primary key,
		managerId bigserial not null references managers,
		title varchar not null,
		description varchar not null,
		files varchar,
		specialityId bigserial,  --references specialities,
		experienceLevelId bigserial,
		paymentAmount decimal(8,2),
		country varchar,
		city varchar,
		jobTypeId bigserial,
		date timestamp,
		status varchar
	);`
	if _, err := db.Exec(jobsQuery); err != nil {
		return err
	}

	specialitiesQuery := `CREATE TABLE IF NOT EXISTS specialities (
		id bigserial not null primary key,
		name varchar
	);`
	if _, err := db.Exec(specialitiesQuery); err != nil {
		return err
	}

	responsesQuery := `CREATE TABLE IF NOT EXISTS responses (
		id bigserial not null primary key,
		freelancerId bigserial not null references freelancers,
		jobId bigserial not null references jobs,
		files varchar,
		date timestamp not null,
		statusManager varchar not null,
		statusFreelancer varchar not null,
		paymentAmount decimal(8,2) not null 
	);`
	if _, err := db.Exec(responsesQuery); err != nil {
		return err
	}

	contractsQuery := `CREATE TABLE IF NOT EXISTS contracts (
		id bigserial not null primary key,
		responseId bigint not null references responses,
		companyId bigint not null references companies,
		freelancerId bigint not null,
		startTime timestamp not null,
		endTime timestamp not null,
		status varchar not null,
		grade int not null,
		paymentAmount decimal(8,2) not null 
	);`
	if _, err := db.Exec(contractsQuery); err != nil {
		return err
	}

	return nil
}

func dropAllTables(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS users;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = `DROP TABLE IF EXISTS managers;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = `DROP TABLE IF EXISTS freelancers;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = `DROP TABLE IF EXISTS jobs;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = `DROP TABLE IF EXISTS companies;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = `DROP TABLE IF EXISTS specialities;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	query = `DROP TABLE IF EXISTS contracts;`
	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}

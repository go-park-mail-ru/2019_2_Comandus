package companyRepository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
)

type CompanyRepository struct {
	db *sql.DB
}


func NewCompanyRepository(db *sql.DB) company.Repository {
	return &CompanyRepository{db}
}

func (r *CompanyRepository) Create(company *model.Company) error {
	return r.db.QueryRow(
		"INSERT INTO companies (companyName, site, tagLine, description, country, city, address, phone) " +
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		company.CompanyName,
		company.Site,
		company.TagLine,
		company.Description,
		company.Country,
		company.City,
		company.Address,
		company.Phone,
	).Scan(&company.ID)
}

func (r *CompanyRepository) Find(id int64) (*model.Company, error) {
	c := &model.Company{}
	if err := r.db.QueryRow(
		"SELECT id, companyName, site, tagLine, description, country, city, address, " +
			"phone FROM companies WHERE id = $1",
		id,
	).Scan(
		&c.ID,
		&c.CompanyName,
		&c.Site,
		&c.TagLine,
		&c.Description,
		&c.Country,
		&c.City,
		&c.Address,
		&c.Phone,
	); err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CompanyRepository) Edit(c * model.Company) error {
	return r.db.QueryRow("UPDATE companies SET companyName = $1, site = $2, tagLine = $3, " +
		"description = $4, country = $5, city = $6, address = $7, phone = $8 WHERE id = $9 RETURNING id",
		c.CompanyName,
		c.Site,
		c.TagLine,
		c.Description,
		c.Country,
		c.City,
		c.Address,
		c.Phone,
		c.ID,
	).Scan(&c.ID)
}

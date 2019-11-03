package sqlstore

import (
	"database/sql"
	_ "github.com/lib/pq"
)

// Store ...
type Store struct {
	db                   *sql.DB
	userRepository       *UserRepository
	freelancerRepository *FreelancerRepository
	managerRepository    *ManagerRepository
	jobRepository 		 *JobRepository
	responseRepository 	 *ResponseRepository
	companyRepository 	 *CompanyRepository
	contractRepository 	 *ContractRepository
	config               *Config
	//Mu                   *sync.Mutex
}


func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Manager() *ManagerRepository {
	if s.managerRepository != nil {
		return s.managerRepository
	}

	s.managerRepository = &ManagerRepository{
		store: s,
	}

	return s.managerRepository
}

func (s *Store) Freelancer() *FreelancerRepository {
	if s.freelancerRepository != nil {
		return s.freelancerRepository
	}

	s.freelancerRepository = &FreelancerRepository{
		store: s,
	}

	return s.freelancerRepository
}

func (s *Store) Job() *JobRepository {
	if s.jobRepository != nil {
		return s.jobRepository
	}

	s.jobRepository = &JobRepository{
		store: s,
	}

	return s.jobRepository
}

func (s *Store) Response() *ResponseRepository {
	if s.responseRepository != nil {
		return s.responseRepository
	}

	s.responseRepository = &ResponseRepository{
		store: s,
	}

	return s.responseRepository
}

func (s *Store) Company() *CompanyRepository {
	if s.companyRepository != nil {
		return s.companyRepository
	}

	s.companyRepository = &CompanyRepository{
		store: s,
	}

	return s.companyRepository
}

func (s *Store) Contract() *ContractRepository {
	if s.contractRepository != nil {
		return s.contractRepository
	}

	s.contractRepository = &ContractRepository {
		store: s,
	}

	return s.contractRepository
}
package sqlstore

import (
	"database/sql"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store"
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

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Manager() store.ManagerRepository {
	if s.managerRepository != nil {
		return s.managerRepository
	}

	s.managerRepository = &ManagerRepository{
		store: s,
	}

	return s.managerRepository
}

func (s *Store) Freelancer() store.FreelancerRepository {
	if s.freelancerRepository != nil {
		return s.freelancerRepository
	}

	s.freelancerRepository = &FreelancerRepository{
		store: s,
	}

	return s.freelancerRepository
}

func (s *Store) Job() store.JobRepository {
	if s.jobRepository != nil {
		return s.jobRepository
	}

	s.jobRepository = &JobRepository{
		store: s,
	}

	return s.jobRepository
}

func (s *Store) Response() store.ResponseRepository {
	if s.responseRepository != nil {
		return s.responseRepository
	}

	s.responseRepository = &ResponseRepository{
		store: s,
	}

	return s.responseRepository
}

func (s *Store) Company() store.CompanyRepository {
	if s.companyRepository != nil {
		return s.companyRepository
	}

	s.companyRepository = &CompanyRepository{
		store: s,
	}

	return s.companyRepository
}

func (s *Store) Contract() store.ContractRepository {
	if s.contractRepository != nil {
		return s.contractRepository
	}

	s.contractRepository = &ContractRepository {
		store: s,
	}

	return s.contractRepository
}
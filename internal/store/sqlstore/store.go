package sqlstore

import (
	"database/sql"
	_ "github.com/lib/pq"
	"sync"
)

// Store ...
type Store struct {
	db                   *sql.DB
	userRepository       *UserRepository
	freelancerRepository *FreelancerRepository
	managerRepository    *ManagerRepository
	jobRepository 		 *JobRepository
	config               *Config
	Mu                   *sync.Mutex
}


func New(db *sql.DB) *Store {
	return &Store{
		db: db,
		Mu: new(sync.Mutex),
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
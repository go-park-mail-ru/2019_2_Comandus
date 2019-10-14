package sqlstore

import (
	"database/sql"
	_ "github.com/lib/pq"
)

// Store ...
type Store struct {
	db             *sql.DB
	userRepository *UserRepository
	freelancerRepository *FreelancerRepository
	managerRepository *ManagerRepository
	config *Config
}


func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s * Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s * Store) Close() {
	s.db.Close()
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
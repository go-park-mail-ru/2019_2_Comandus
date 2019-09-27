package store


import (
	"database/sql"
	_ "github.com/lib/pq"
)

// Store ...
type Store struct {
	db             *sql.DB
	//userRepository *UserRepository
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
// New ...
/*func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User ...
func (s *Store) User() UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}*/
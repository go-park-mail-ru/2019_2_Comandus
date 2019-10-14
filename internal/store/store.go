package store

type Store interface {
	User() UserRepository
	Manager() ManagerRepository
	Freelancer() FreelancerRepository
}

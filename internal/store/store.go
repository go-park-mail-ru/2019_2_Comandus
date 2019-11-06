package store

type Store interface {
	User() UserRepository
	Manager() ManagerRepository
	Freelancer() FreelancerRepository
	Job() JobRepository
	Response() ResponseRepository
	Company() CompanyRepository
	Contract() ContractRepository
}

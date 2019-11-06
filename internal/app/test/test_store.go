package test

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/store"
	"github.com/golang/mock/gomock"
	"testing"
)

type Store struct {
	userRepository       *MockUserRepository
	freelancerRepository *MockFreelancerRepository
	managerRepository    *MockManagerRepository
	jobRepository 		 *MockJobRepository
	responseRepository 	 *MockResponseRepository
	companyRepository 	 *MockCompanyRepository
	contractRepository 	 *MockContractRepository
}


func New(t *testing.T) *Store {
	t.Helper()
	ctrl := gomock.NewController(t)
	return &Store{
		userRepository:NewMockUserRepository(ctrl),
		freelancerRepository:NewMockFreelancerRepository(ctrl),
		managerRepository:NewMockManagerRepository(ctrl),
		jobRepository:NewMockJobRepository(ctrl),
		responseRepository:NewMockResponseRepository(ctrl),
		companyRepository:NewMockCompanyRepository(ctrl),
		contractRepository:NewMockContractRepository(ctrl),
	}
}

func (s *Store) User() store.UserRepository {
	return s.userRepository
}

func (s *Store) Manager() store.ManagerRepository {
	return s.managerRepository
}

func (s *Store) Freelancer() store.FreelancerRepository {
	return s.freelancerRepository
}

func (s *Store) Job() store.JobRepository {
	return s.jobRepository
}

func (s *Store) Response() store.ResponseRepository {
	return s.responseRepository
}

func (s *Store) Company() store.CompanyRepository {
	return s.companyRepository
}

func (s *Store) Contract() store.ContractRepository {
	return s.contractRepository
}


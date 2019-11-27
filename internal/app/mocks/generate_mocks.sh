echo "generating repository mocks..."
mockgen -source=internal/app/user/repository.go -package=repository_mocks -mock_names=Repository=MockUserRepository > internal/app/mocks/repository_mocks/user_rep_mock.go
mockgen -source=internal/app/freelancer/repository.go -package=repository_mocks -mock_names=Repository=MockFreelancerRepository  > internal/app/mocks/repository_mocks/freelancer_rep_mock.go
mockgen -source=internal/app/manager/repository.go -package=repository_mocks -mock_names=Repository=MockManagerRepository  > internal/app/mocks/repository_mocks/manager_rep_mock.go
mockgen -source=internal/app/company/repository.go -package=repository_mocks -mock_names=Repository=MockCompanyRepository  > internal/app/mocks/repository_mocks/company_rep_mock.go
mockgen -source=internal/app/user-job/repository.go -package=repository_mocks -mock_names=Repository=MockJobRepository  > internal/app/mocks/repository_mocks/job_rep_mock.go
mockgen -source=internal/app/user-response/repository.go -package=repository_mocks -mock_names=Repository=MockResponseRepository  > internal/app/mocks/repository_mocks/response_rep_mock.go
mockgen -source=internal/app/user-contract/repository.go -package=repository_mocks -mock_names=Repository=MockContractRepository  > internal/app/mocks/repository_mocks/contract_rep_mock.go

echo "generating usecase mocks..."
mockgen -source=internal/app/freelancer/usecase.go -package=ucase_mocks -mock_names=Usecase=MockFreelancerUsecase > internal/app/mocks/ucase_mocks/freelancer_ucase_mock.go
mockgen -source=internal/app/user/usecase.go -package=ucase_mocks -mock_names=Usecase=MockUserUsecase > internal/app/mocks/ucase_mocks/user_ucase_mock.go
mockgen -source=internal/app/user/usecase.go -package=ucase_mocks -mock_names=Usecase=MockUserUsecase > internal/app/mocks/ucase_mocks/user_ucase_mock.go
mockgen -source=internal/app/user-contract/usecase.go -package=ucase_mocks -mock_names=Usecase=MockContractUsecase > internal/app/mocks/ucase_mocks//contract_ucase_mock.go
mockgen -source=internal/app/user-job/usecase.go -package=ucase_mocks -mock_names=Usecase=MockJobUsecase > internal/app/mocks/ucase_mocks/job_ucase_mock.go
mockgen -source=internal/app/user-response/usecase.go -package=ucase_mocks -mock_names=Usecase=MockResponseUsecase > internal/app/mocks/ucase_mocks/response_ucase_mock.go
mockgen -source=internal/app/company/usecase.go -package=ucase_mocks -mock_names=Usecase=MockCompanyUsecase > internal/app/mocks/ucase_mocks/company_ucase_mock.go

echo "generating clients mocks..."
mockgen -source=internal/app/clients/interfaces/managerClient.go -package=client_mocks -mock_names=Clients=MockManagerClient > internal/app/mocks/client_mocks/manager_client_mock.go
mockgen -source=internal/app/clients/interfaces/companyClient.go -package=client_mocks -mock_names=Clients=MockCompanyClient > internal/app/mocks/client_mocks/company_client_mock.go
mockgen -source=internal/app/clients/interfaces/freelancerClient.go -package=client_mocks -mock_names=Clients=MockFreelancerClient > internal/app/mocks/client_mocks/freelancer_client_mock.go
mockgen -source=internal/app/clients/interfaces/jobClient.go -package=client_mocks -mock_names=Clients=MockJobClient > internal/app/mocks/client_mocks/job_client_mock.go
mockgen -source=internal/app/clients/interfaces/responseClient.go -package=client_mocks -mock_names=Clients=MockResponseClient > internal/app/mocks/client_mocks/response_client_mock.go
mockgen -source=internal/app/clients/interfaces/userClient.go -package=client_mocks -mock_names=Clients=MockUserClient > internal/app/mocks/client_mocks/user_client_mock.go
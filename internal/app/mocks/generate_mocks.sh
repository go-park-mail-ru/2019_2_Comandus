echo "generating repository mocks..."
mockgen -source=internal/app/user/repository.go -package=mocks -mock_names=Repository=MockUserRepository > internal/app/mocks/user_rep_mock.go
mockgen -source=internal/app/freelancer/repository.go -package=mocks -mock_names=Repository=MockFreelancerRepository  > internal/app/mocks/freelancer_rep_mock.go
mockgen -source=internal/app/manager/repository.go -package=mocks -mock_names=Repository=MockManagerRepository  > internal/app/mocks/manager_rep_mock.go
mockgen -source=internal/app/company/repository.go -package=mocks -mock_names=Repository=MockCompanyRepository  > internal/app/mocks/company_rep_mock.go
mockgen -source=internal/app/user-job/repository.go -package=mocks -mock_names=Repository=MockJobRepository  > internal/app/mocks/job_rep_mock.go
mockgen -source=internal/app/user-response/repository.go -package=mocks -mock_names=Repository=MockResponseRepository  > internal/app/mocks/response_rep_mock.go
mockgen -source=internal/app/user-contract/repository.go -package=mocks -mock_names=Repository=MockContractRepository  > internal/app/mocks/contract_rep_mock.go

echo "generating usecase mocks..."
mockgen -source=internal/app/freelancer/usecase.go -package=test -mock_names=Usecase=MockFreelancerUsecase > internal/app/freelancer/delivery/http/test/freelancer_ucase_mock.go
mockgen -source=internal/app/user/usecase.go -package=test -mock_names=Usecase=MockUserUsecase > internal/app/user/delivery/http/test/user_ucase_mock.go
mockgen -source=internal/app/user/usecase.go -package=test -mock_names=Usecase=MockUserUsecase > internal/app/general/delivery/http/test/user_ucase_mock.go
mockgen -source=internal/app/user-contract/usecase.go -package=test -mock_names=Usecase=MockContractUsecase > internal/app/user-contract/delivery/http/test/contract_ucase_mock.go
mockgen -source=internal/app/user-job/usecase.go -package=test -mock_names=Usecase=MockJobUsecase > internal/app/user-job/delivery/http/test/job_ucase_mock.go
mockgen -source=internal/app/user-response/usecase.go -package=test -mock_names=Usecase=MockResponseUsecase > internal/app/user-response/delivery/http/test/response_ucase_mock.go
mockgen -source=internal/app/company/usecase.go -package=test -mock_names=Usecase=MockCompanyUsecase > internal/app/company/delivery/http/test/company_ucase_mock.go

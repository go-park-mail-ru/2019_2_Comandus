#### Генерация mock файлов для структур repository
`$:mockgen -source=internal/app/user/repository.go -package=test -mock_names=Repository=MockUserRepository > user_rep_mock.go`

`$:mockgen -source=internal/app/freelancer/repository.go -package=test -mock_names=Repository=MockFreelancerRepository  > freelancer_rep_mock.go`

`$:mockgen -source=internal/app/manager/repository.go -package=test -mock_names=Repository=MockManagerRepository  > manager_rep_mock.go`

`$:mockgen -source=internal/app/company/repository.go -package=test -mock_names=Repository=MockCompanyRepository  > company_rep_mock.go`

`$:mockgen -source=internal/app/user-job/repository.go -package=test -mock_names=Repository=MockJobRepository  > job_rep_mock.go`

`$:mockgen -source=internal/app/user-response/repository.go -package=test -mock_names=Repository=MockResponseRepository  > response_rep_mock.go`

`$:mockgen -source=internal/app/user-contract/repository.go -package=test -mock_names=Repository=MockContractRepository  > contract_rep_mock.go`

#### Генерация mock файлов для структур usecase
`$:mockgen -source=internal/app/freelancer/usecase.go -package=test -mock_names=Usecase=MockFreelancerUsecase > internal/app/freelancer/delivery/http/test/freelancer_ucase_mock.go` 

`$:mockgen -source=internal/app/user/usecase.go -package=test -mock_names=Usecase=MockUserUsecase > internal/app/user/delivery/http/test/user_ucase_mock.go`

`$:mockgen -source=internal/app/user/usecase.go -package=test -mock_names=Usecase=MockUserUsecase > internal/app/general/delivery/http/test/user_ucase_mock.go`

`$:mockgen -source=internal/app/user-contract/usecase.go -package=test -mock_names=Usecase=MockContractUsecase > internal/app/user-contract/delivery/http/test/contract_ucase_mock.go` 

`$:mockgen -source=internal/app/user-job/usecase.go -package=test -mock_names=Usecase=MockJobUsecase > internal/app/user-job/delivery/http/test/job_ucase_mock.go`

`$:mockgen -source=internal/app/user-response/usecase.go -package=test -mock_names=Usecase=MockResponseUsecase > internal/app/user-response/delivery/http/test/response_ucase_mock.go` 

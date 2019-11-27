package companyUsecase

import (
	server_clients "github.com/go-park-mail-ru/2019_2_Comandus/internal/app/clients/server-clients"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/company"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/manager/delivery/grpc/manager_grpc"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/mocks/client_mocks"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/app/mocks/repository_mocks"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testClients(t *testing.T) {
	t.Helper()
	ctrl := gomock.NewController(t)
	serverclients := server_clients.ServerClients{
		AuthClient:       client_mocks.NewMockAuthUser(ctrl),
		CompanyClient:    nil,
		FreelancerClient: nil,
		JobClient:        nil,
		ManagerClient:    nil,
		ResponseClient:   nil,
		UserClient:       nil,
	}
}

func testUcase(t *testing.T) (*repository_mocks.MockCompanyRepository, *client_mocks.MockManagerClient ,company.Usecase){
	t.Helper()
	ctrl := gomock.NewController(t)
	companyRep := repository_mocks.NewMockCompanyRepository(ctrl)
	managerClient := client_mocks.NewMockManagerClient(ctrl)
	companyUcase := NewCompanyUsecase(companyRep, managerClient)
	return companyRep, managerClient,  companyUcase
}

func TestCompanyUsecase_Create(t *testing.T) {
	companyRep, _ ,companyUcase := testUcase(t)

	testCases := []struct {
		name			string
		newCompany		*model.Company
		userType		string
		expectError		error
	}{
		{
			name:			"valid",
			newCompany:		&model.Company{},
			userType:		model.UserCustomer,
			expectError:	nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			companyRep.
				EXPECT().
				Create(tc.newCompany).
				Do(func(arg *model.Company){
					arg.ID = 1
			}).
				Return(tc.expectError)

			c, err := companyUcase.Create()

			if tc.expectError == nil {
				assert.Equal(t, int64(1) , c.ID)
				assert.Equal(t, nil, err)
				return
			}

			if err != nil {
				assert.Equal(t, tc.expectError.Error(), err.Error())
				return
			}

			t.Fatal()
		})
	}
}

func TestCompanyUsecase_Edit(t *testing.T) {
	companyRep, managerClient, companyUcase := testUcase(t)

	user := &model.User{
		ID:               1,
		FirstName:        "ddd",
		Email:            "ddd@hj.cv",
	}

	testCases := []struct {
		name			string
		newCompany		*model.Company
		userType		string
		expectError		error
		expectedManager *manager_grpc.Manager
	}{
		{
			name:			"valid",
			newCompany:		&model.Company{
				ID: 		1,
				CompanyName: "new name",
				Site:        "www.new-site.ru",
			},
			expectedManager: &manager_grpc.Manager{
				ID:                   1,
				CompanyId:            1,
			},
			userType:		model.UserCustomer,
			expectError:	nil,
		},
		{
			name:			"invalid: user is freelancer",
			newCompany:		&model.Company{
				ID: 		 1,
				CompanyName: "new name",
				Site:        "www.new-site.ru",
			},
			expectedManager: &manager_grpc.Manager{
				ID:                   1,
				CompanyId:            1,
			},
			userType:		model.UserFreelancer,
			expectError:	errors.New(" only manager can edit company"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user.UserType = tc.userType
			managerClient.
				EXPECT().
				GetManagerByUserFromServer(user.ID).
				Return(tc.expectedManager, nil)


			companyRep.
				EXPECT().
				Edit(tc.newCompany).
				Return(tc.expectError)

			err := companyUcase.Edit(user, tc.newCompany)

			if tc.expectError == nil {
				assert.Equal(t, nil, err)
				return
			}

			expectError := "HandleEditCompany<-Edit: :  only manager can edit company"

			if err != nil {
				assert.Equal(t, expectError, err.Error())
				return
			}
			t.Fatal()
		})
	}
}

func TestCompanyUsecase_Find(t *testing.T) {
	companyRep, _, companyUcase := testUcase(t)

	testCases := []struct {
		name        string
		company     *model.Company
		expectError error
	}{
		{
			name:			"valid",
			company:		&model.Company{
				ID:          1,
				CompanyName: "test company",
				Site:        "test",
			},
			expectError:	nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			companyRep.
				EXPECT().
				Find(tc.company.ID).
				Return(tc.company, err)

			company, err := companyUcase.Find(tc.company.ID)

			assert.Equal(t, nil, err)
			assert.Equal(t, company, testCases[0].company)
		})
	}
}

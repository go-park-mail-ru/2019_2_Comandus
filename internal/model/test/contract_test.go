package test

import (
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testContract(t *testing.T) *model.Contract {
	t.Helper()
	return &model.Contract{
		ID:           0,
		ResponseID:   1,
		CompanyID:    1,
		FreelancerID: 1,
		Status:       "in work",
		Grade:        0,
	}
}

func  TestContract_Validate(t *testing.T) {
	testCases := []struct {
		name     string
		contract func() *model.Contract
		isValid  bool
	}{
		{
			name: "valid",
			contract: func() *model.Contract {
				return testContract(t)
			},
			isValid: true,
		},
		{
			name: "empty status",
			contract: func() *model.Contract {
				c := testContract(t)
				c.Status = ""
				return c
			},
			isValid: false,
		},
		{
			name: "invalid references",
			contract: func() *model.Contract {
				c := testContract(t)
				c.ResponseID = 0
				return c
			},
			isValid: false,
		},
		{
			name: "incorrect last id",
			contract: func() *model.Contract {
				c := testContract(t)
				c.ID = 100
				return c
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.contract().Validate(0))
			} else {
				assert.Error(t, tc.contract().Validate(0))
			}
		})
	}
}


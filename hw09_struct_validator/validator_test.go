package hw09structvalidator

import (
	"fmt"
	"testing"

	errors2 "github.com/IASamoylov/otus_home_work/hw09_struct_validator/errors"
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string      `json:"id" validate:"len:36"`
		Name   interface{} `validate:"nested"`
		Age    int         `validate:"min:18|max:50"`
		Email  string      `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole    `validate:"in:admin,stuff"`
		Phones []string    `validate:"len:11"`
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr string
	}{
		{
			in: Response{
				Code: 200,
			},
			expectedErr: "tag `validate:\"in:200,404,500\"` configured incorrectly validation rule for field Code",
		},
		{
			in: App{
				Version: "release-1.0.0",
			},
			expectedErr: "Version: string length must be equal 5",
		},
		{
			in: User{
				ID: "0dc6bb8d-e8e7-4356-a02a-eff220e769f9",
				Name: struct {
					FirstName string `validate:"min:3|max:50"`
					LastName  string `validate:"min:3|max:50"`
				}{FirstName: "ad", LastName: "sad"},
			},
			expectedErr: "tag `validate:\"min:3|max:50\"` not supported for this type string",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.EqualError(t, err, tt.expectedErr)
		})
	}
}

func Test(t *testing.T) {
	user := User{
		ID: "0dc6bb8d-e8e7-4356-a02a-eff220e769f9",
		Name: struct {
			FirstName string `validate:"regexp:^\\S{3,}$"`
			LastName  string `validate:"regexp:^\\S{3,}$"`
		}{FirstName: "FirstName", LastName: "La"},
		Age:   13,
		Email: "ia.sa@otu.com",
		Role:  "admin",
		Phones: []string{
			"89006731166",
			"8006731166",
		},
	}
	user1 := &user

	err := Validate(user)
	var validationErrors *errors2.ValidationErrors
	require.ErrorAs(t, err, &validationErrors)
	require.EqualError(t, validationErrors.Errors[0], "Name.LastName: value does not match the mask ^\\S{3,}$")
	require.EqualError(t, validationErrors.Errors[1], "Age: value must be greater or equal 18")
	require.EqualError(t, validationErrors.Errors[2], "Email: value does not match the mask ^\\w+@\\w+\\.\\w+$")
	require.EqualError(t, validationErrors.Errors[3], "Phones[1]: string length must be equal 11")

	err = Validate(user1)
	require.ErrorAs(t, err, &validationErrors)
	require.EqualError(t, validationErrors.Errors[0], "Name.LastName: value does not match the mask ^\\S{3,}$")
	require.EqualError(t, validationErrors.Errors[1], "Age: value must be greater or equal 18")
	require.EqualError(t, validationErrors.Errors[2], "Email: value does not match the mask ^\\w+@\\w+\\.\\w+$")
	require.EqualError(t, validationErrors.Errors[3], "Phones[1]: string length must be equal 11")
}

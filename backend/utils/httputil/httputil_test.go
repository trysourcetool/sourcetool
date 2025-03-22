package httputil

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

type testPasswordStruct struct {
	Password string `validate:"required,password"`
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password with minimum requirements",
			password: "Password1",
			wantErr:  false,
		},
		{
			name:     "valid password with special characters",
			password: "Password1!@#$",
			wantErr:  false,
		},
		{
			name:     "valid password with all allowed special characters",
			password: "Pass1!?_+*'\"` $%&-^@;:,./=~|[](){}<>",
			wantErr:  true,
		},
		{
			name:     "valid password with some allowed special characters",
			password: "Pass1!?_+*'\"@;:,./=~|[](){}<>",
			wantErr:  false,
		},
		{
			name:     "too short password",
			password: "Pass1",
			wantErr:  true,
		},
		{
			name:     "password without letters",
			password: "12345678",
			wantErr:  true,
		},
		{
			name:     "password without numbers",
			password: "Password",
			wantErr:  true,
		},
		{
			name:     "password with invalid special characters",
			password: "Password1±§",
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "password with spaces",
			password: "Password 1",
			wantErr:  true,
		},
		{
			name:     "password with japanese characters",
			password: "Password1あいう",
			wantErr:  true,
		},
	}

	v := validator.New()
	v.RegisterValidation("password", validatePassword)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := &testPasswordStruct{
				Password: tt.password,
			}
			err := v.Struct(test)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePassword() error = %v, wantErr %v", err, tt.wantErr)
				if err != nil {
					t.Logf("validation error details: %v", err)
				}
			}
		})
	}
}

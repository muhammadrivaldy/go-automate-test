package tests

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestGetUsers(t *testing.T) {

	type useCase struct {
		Name             string
		ExpectedCode     int
		ExpectedResponse string
	}

	useCases := []useCase{
		{
			Name:             "Success get users",
			ExpectedCode:     http.StatusOK,
			ExpectedResponse: `{"users":[{"id":1,"name":"Anton","address":"Jl. Padepokan"},{"id":2,"name":"Steven","address":"Jl. Studio Raya"}]}`,
		},
	}

	for _, useCase := range useCases {

		t.Run(useCase.Name, func(t *testing.T) {

			baseURL := os.Getenv("BASE_URL")

			request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/users", baseURL), nil)
			if err != nil {
				t.Error(err)
				return
			}

			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				t.Error(err)
				return
			}

			defer response.Body.Close()

			responseInByte, err := io.ReadAll(response.Body)
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, useCase.ExpectedCode, response.StatusCode)
			assert.Equal(t, useCase.ExpectedResponse, string(responseInByte))
		})
	}
}

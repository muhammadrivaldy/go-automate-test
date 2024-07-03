package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestCreateUsers(t *testing.T) {

	type user struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}

	type useCase struct {
		Name         string
		Parameters   user
		ExpectedCode int
	}

	useCases := []useCase{
		{
			Name: "Success create user - Anton",
			Parameters: user{
				Name:    "Anton",
				Address: "Jl. Padepokan",
			},
			ExpectedCode: http.StatusCreated,
		},
		{
			Name: "Success create user - Steven",
			Parameters: user{
				Name:    "Steven",
				Address: "Jl. Studio Raya",
			},
			ExpectedCode: http.StatusOK,
		},
	}

	for _, useCase := range useCases {

		t.Run(useCase.Name, func(t *testing.T) {

			baseURL := os.Getenv("BASE_URL")

			paramInBytes, err := json.Marshal(useCase.Parameters)
			if err != nil {
				t.Error(err)
				return
			}

			request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users", baseURL), bytes.NewReader(paramInBytes))
			if err != nil {
				t.Error(err)
				return
			}

			request.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				t.Error(err)
				return
			}

			defer response.Body.Close()

			assert.Equal(t, http.StatusCreated, response.StatusCode)
		})
	}
}

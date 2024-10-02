package main_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/raviand/test-project/pkg"
	"github.com/stretchr/testify/require"
)

type ImportantService interface {
	Do(int, string) (string, error)
}

type Dummy struct {
}

func (d *Dummy) Do(int, string) (string, error) {
	return "dummy", nil
}

func TestIntegration(t *testing.T) {
	t.Run("test dummy", func(t *testing.T) {
		testId := "123e4567-e89b-12d3-a456-426614174020"
		u := pkg.User{
			Id:         testId,
			Firstname:  "John",
			Lastname:   "Doe",
			Username:   "johndoe",
			Password:   "P@ssw0rd!",
			Email:      "johndoe@example.com",
			IP:         "192.168.1.1",
			MacAddress: "00:14:22:01:23:45",
			Website:    "https://johndoe.com",
			Image:      "https://example.com/images/johndoe.jpg",
		}
		b, err := json.Marshal(u)
		require.Nil(t, err)
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/user", strings.NewReader(string(b)))
		require.Nil(t, err)
		req.Header.Set("Authorization", "this-is-my-token")
		req.Header.Set("X-User", "johndoe")
		res, err := http.DefaultClient.Do(req)
		require.Nil(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)
		t.Log("response status code:", res.StatusCode)
		req, err = http.NewRequest(http.MethodGet, "http://localhost:8080/user/"+testId, nil)
		require.Nil(t, err)
		req.Header.Set("Authorization", "this-is-my-token")
		req.Header.Set("X-User", "johndoe")
		res, err = http.DefaultClient.Do(req)
		require.Nil(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
		req, err = http.NewRequest(http.MethodDelete, "http://localhost:8080/user/"+testId, nil)
		require.Nil(t, err)
		req.Header.Set("Authorization", "this-is-my-token")
		req.Header.Set("X-User", "johndoe")
		res, err = http.DefaultClient.Do(req)
		require.Nil(t, err)
		require.Equal(t, http.StatusNoContent, res.StatusCode)
		req, err = http.NewRequest(http.MethodGet, "http://localhost:8080/user/"+testId, nil)
		require.Nil(t, err)
		req.Header.Set("Authorization", "this-is-my-token")
		req.Header.Set("X-User", "johndoe")
		res, err = http.DefaultClient.Do(req)
		require.Nil(t, err)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}

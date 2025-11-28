package controller_test

import (
	"fmt"
	"grf/core/pagination"
	"grf/core/tests"
	authdto "grf/domain/auth/dto"
	"net/http"
	"testing"

	"github.com/goccy/go-json"
)

func TestUserCRUD(t *testing.T) {
	clearAuthTables(testApp.DB)
	_, err := createTestFixtures(testApp.DB)
	if err != nil {
		t.Fatalf("Falha ao criar fixtures: %v", err)
	}

	adminToken, _ := loginAs(t, "admin", "admin123")
	userToken, _ := loginAs(t, "user", "user123")

	var createdUserID uint64

	t.Run("POST /users (Admin 201)", func(t *testing.T) {
		dto := authdto.UserCreateDTO{
			Username: "newuser", Email: "new@user.com", Password: "password123",
		}
		resp, body := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodPost, URL: "/v1/users", Token: adminToken, Body: dto,
		})
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Esperado 201, obteve %d: %s", resp.StatusCode, body)
		}

		var respDTO authdto.UserResponseDTO
		err := json.Unmarshal([]byte(body), &respDTO)
		if err != nil {
			return
		}
		if respDTO.Username != "newuser" {
			t.Errorf("Esperado 'newuser', obteve %s", respDTO.Username)
		}
		createdUserID = respDTO.ID
	})

	t.Run("POST /users (User 403)", func(t *testing.T) {
		dto := authdto.UserCreateDTO{Username: "failuser", Email: "fail@user.com", Password: "password123"}
		resp, _ := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodPost, URL: "/v1/users", Token: userToken, Body: dto,
		})
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Esperado 403, obteve %d", resp.StatusCode)
		}
	})

	t.Run("GET /users (Admin 200)", func(t *testing.T) {
		resp, body := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodGet, URL: "/v1/users", Token: adminToken,
		})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Esperado 200, obteve %d: %s", resp.StatusCode, body)
		}

		var listResp pagination.Response[authdto.UserResponseDTO]
		err := json.Unmarshal([]byte(body), &listResp)
		if err != nil {
			return
		}
		if *listResp.Count < 3 {
			t.Errorf("Esperado 3+ usuários, obteve %d", *listResp.Count)
		}
	})

	t.Run("GET /users/:id (Admin 200)", func(t *testing.T) {
		url := fmt.Sprintf("/v1/users/%d", createdUserID)
		resp, _ := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodGet, URL: url, Token: adminToken,
		})
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Esperado 200, obteve %d", resp.StatusCode)
		}
	})

	t.Run("PATCH /users/:id (Admin 200)", func(t *testing.T) {
		url := fmt.Sprintf("/v1/users/%d", createdUserID)
		patchDTO := map[string]string{"first_name": "Patched"}

		resp, body := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodPatch, URL: url, Token: adminToken, Body: patchDTO,
		})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Esperado 200, obteve %d: %s", resp.StatusCode, body)
		}

		var respDTO authdto.UserResponseDTO
		err := json.Unmarshal([]byte(body), &respDTO)
		if err != nil {
			return
		}
		if respDTO.FirstName != "Patched" {
			t.Errorf("Patch falhou")
		}
	})

	t.Run("DELETE /users/:id (Admin 204)", func(t *testing.T) {
		url := fmt.Sprintf("/v1/users/%d", createdUserID)
		resp, _ := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodDelete, URL: url, Token: adminToken,
		})
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("Esperado 204, obteve %d", resp.StatusCode)
		}
	})

	t.Run("GET /users/:id (Admin 404 Afer Delete)", func(t *testing.T) {
		url := fmt.Sprintf("/v1/users/%d", createdUserID)
		resp, _ := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodGet, URL: url, Token: adminToken,
		})
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Esperado 404, obteve %d", resp.StatusCode)
		}
	})
}

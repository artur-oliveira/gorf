package controller_test

import (
	"encoding/json"
	"fmt"
	"grf/core/tests"
	authdto "grf/domain/auth/dto"
	"net/http"
	"testing"
)

func TestGroupCRUD(t *testing.T) {
	clearAuthTables(testApp.DB)
	fixtures, err := createTestFixtures(testApp.DB)
	if err != nil {
		t.Fatalf("Falha ao criar fixtures: %v", err)
	}

	adminToken, _ := loginAs(t, "admin", "admin123")
	userToken, _ := loginAs(t, "user", "user123")

	var createdGroupID uint64
	permID := fixtures.PermListUser.ID

	t.Run("POST /groups (Admin 201 M2M)", func(t *testing.T) {
		dto := authdto.GroupCreateDTO{
			Name:          "Grupo M2M",
			PermissionIDs: []uint64{permID},
		}
		resp, body := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodPost, URL: "/v1/groups", Token: adminToken, Body: dto,
		})
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Esperado 201, obteve %d: %s", resp.StatusCode, body)
		}

		var respDTO authdto.GroupResponseDTO
		err := json.Unmarshal([]byte(body), &respDTO)
		if err != nil {
			return
		}
		if respDTO.Name != "Grupo M2M" {
			t.Errorf("Nome incorreto")
		}
		if len(respDTO.Permissions) != 1 || respDTO.Permissions[0].ID != permID {
			t.Fatal("Associação M2M falhou na criação")
		}
		createdGroupID = respDTO.ID
	})

	t.Run("POST /groups (User 403)", func(t *testing.T) {
		dto := authdto.GroupCreateDTO{Name: "Grupo Falho"}
		resp, _ := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodPost, URL: "/v1/groups", Token: userToken, Body: dto,
		})
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Esperado 403, obteve %d", resp.StatusCode)
		}
	})

	t.Run("PUT /groups/:id (Admin 200 M2M Replace)", func(t *testing.T) {
		url := fmt.Sprintf("/v1/groups/%d", createdGroupID)
		newPermID := fixtures.PermAddUser.ID
		dto := authdto.GroupUpdateDTO{
			Name:          "Grupo M2M Atualizado",
			PermissionIDs: []uint64{newPermID},
		}

		resp, body := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodPut, URL: url, Token: adminToken, Body: dto,
		})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Esperado 200, obteve %d: %s", resp.StatusCode, body)
		}

		var respDTO authdto.GroupResponseDTO
		err := json.Unmarshal([]byte(body), &respDTO)
		if err != nil {
			return
		}
		if respDTO.Name != "Grupo M2M Atualizado" {
			t.Errorf("PUT falhou no nome")
		}
		if len(respDTO.Permissions) != 1 || respDTO.Permissions[0].ID != newPermID {
			t.Fatal("Substituição M2M (PUT) falhou")
		}
	})

	t.Run("DELETE /groups/:id (Admin 204)", func(t *testing.T) {
		url := fmt.Sprintf("/v1/groups/%d", createdGroupID)
		resp, _ := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodDelete, URL: url, Token: adminToken,
		})
		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("Esperado 204, obteve %d", resp.StatusCode)
		}
	})
}

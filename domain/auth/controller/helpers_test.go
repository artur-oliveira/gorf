package controller_test

import (
	"encoding/json"
	"grf/core/models"
	tests2 "grf/core/tests"
	"grf/domain/auth/dto"
	"grf/domain/auth/model"
	"net/http"
	"testing"

	"gorm.io/gorm"
)

var authTables = []string{
	"auth_user_permissions",
	"auth_user_groups",
	"auth_group_permissions",
	"auth_user",
	"auth_group",
}

func clearAuthTables(db *gorm.DB) {
	tests2.ClearTables(db, authTables)
}

type TestFixtures struct {
	AdminUser    *model.User
	NormalUser   *model.User
	PermListUser *model.Permission
	PermAddUser  *model.Permission
}

func getPerm(db *gorm.DB, module string, action string) *model.Permission {
	var perm *model.Permission
	if err := db.Where("module = ? and action = ?", module, action).First(&perm).Error; err != nil {
		panic(err)
	}
	return perm
}

func createTestFixtures(db *gorm.DB) (*TestFixtures, error) {
	perms := []*model.Permission{
		getPerm(db, "user", models.ListAction),
		getPerm(db, "user", models.DetailAction),
		getPerm(db, "user", models.CreateAction),
		getPerm(db, "user", models.UpdateAction),
		getPerm(db, "user", models.PartialUpdateAction),
		getPerm(db, "user", models.DeleteAction),

		getPerm(db, "group", models.ListAction),
		getPerm(db, "group", models.DetailAction),
		getPerm(db, "group", models.CreateAction),
		getPerm(db, "group", models.UpdateAction),
		getPerm(db, "group", models.PartialUpdateAction),
		getPerm(db, "group", models.DeleteAction),

		getPerm(db, "permission", models.ListAction),
		getPerm(db, "permission", models.DetailAction),
		getPerm(db, "permission", models.CreateAction),
		getPerm(db, "permission", models.UpdateAction),
		getPerm(db, "permission", models.PartialUpdateAction),
		getPerm(db, "permission", models.DeleteAction),
	}
	adminGroup := model.Group{Name: "Admin"}
	if err := db.Create(&adminGroup).Error; err != nil {
		return nil, err
	}
	err := db.Model(&adminGroup).Association("Permissions").Append(perms)
	if err != nil {
		return nil, err
	}

	adminUser := model.User{Username: "admin", Email: "admin@test.com", IsActive: true, IsSuperuser: true}
	err = adminUser.SetPassword("admin123")
	if err != nil {
		return nil, err
	}
	if err := db.Create(&adminUser).Error; err != nil {
		return nil, err
	}

	normalUser := model.User{Username: "user", Email: "user@test.com", IsActive: true}
	err = normalUser.SetPassword("user123")
	if err != nil {
		return nil, err
	}
	if err := db.Create(&normalUser).Error; err != nil {
		return nil, err
	}

	return &TestFixtures{
		AdminUser:  &adminUser,
		NormalUser: &normalUser,

		PermListUser: perms[0],
		PermAddUser:  perms[2],
	}, nil
}

func loginAs(t *testing.T, username, password string) (accessToken, refreshToken string) {
	loginDTO := dto.ObtainTokenDTO{
		Login:    username,
		Password: password,
	}

	resp, body := tests2.MakeRequest(t, testApp.FiberApp, tests2.RequestOptions{
		Method: http.MethodPost,
		URL:    "/v1/auth/token",
		Body:   loginDTO,
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Falha ao logar como %s, status %d: %s", username, resp.StatusCode, body)
	}

	var tokenResp dto.TokenResponseDTO
	err := json.Unmarshal([]byte(body), &tokenResp)
	if err != nil {
		return "", ""
	}
	if tokenResp.AccessToken == "" || tokenResp.RefreshToken == "" {
		t.Fatal("Falha ao logar, tokens vazios")
	}

	return tokenResp.AccessToken, tokenResp.RefreshToken
}

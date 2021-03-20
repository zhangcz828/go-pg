package postgres

import (
	"github.com/gin-gonic/gin"
	"go-pg/modules"
	"net/http"
	"net/http/httptest"
	"testing"
)

type DbStoreMock struct{}

func (s DbStoreMock) GetAllHeros() modules.Heros {
	heros := modules.Heros{
		{
			Name:"杨过",
			Detail:"武功：黯然销魂掌, 蛤蟆功 xxxxxxxxxxxxx",
			AttackPower:30,
			DefensePower:15,
			Blood:100,
		},
		{
			Name: "xx",
			Detail: "xx",
			AttackPower: 2,
			DefensePower: 4,
			Blood: 100,
		},
	}

	return heros
}

func TestGetHerosHandler(t *testing.T) {
	// Create Mock Store
	s := DbStoreMock{}

	//// Gin settings.
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/heros", GetHerosHandler(s))

	tests := map[string]struct{
		method string
		path string
		wantCode int
		wantBody string
	} {
		"GetMethod": {
			method: "GET",
			path: "/heros",
			wantCode: http.StatusOK,
			wantBody: `[{"name":"杨过","Detail":"武功：黯然销魂掌, 蛤蟆功 xxxxxxxxxxxxx","AttackPower":30,"DefensePower":15,"Blood":100},{"name":"xx","Detail":"xx","AttackPower":2,"DefensePower":4,"Blood":100}]`,
		},
		"InvalidPath": {
			method: "GET",
			path: "/herosssssss",
			wantCode: http.StatusNotFound,
			wantBody: "404 page not found",
		},
		"InvalidMethod": {
			method: "PUT",
			path: "/heros",
			wantCode: http.StatusNotFound,
			wantBody: "404 page not found",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// 1. httptest 请求和返回
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rr := httptest.NewRecorder()

			// 2. 注册
			router.ServeHTTP(rr, req)

			// 3. 校验判断
			if status := rr.Code; status != tc.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.wantCode)
			}

			if rr.Body.String() != tc.wantBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.wantBody)
			}
		})
	}














}

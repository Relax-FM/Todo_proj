package handler

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/Relax-FM/todo-app-go"
	"github.com/Relax-FM/todo-app-go/pkg/service"
	mock_service "github.com/Relax-FM/todo-app-go/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestHandler_signup(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user todo.User)

	testTable := []struct{
		name 				string
		inputBody 			string
		inputUser 			todo.User
		mockBehavior 		mockBehavior
		expectedStatusCode 	int
		expectedRequestBody string
	} {
		{
			name: "return OK",
			inputBody: `{"name":"Test","username":"test","password":"qwerty"}`,
			inputUser: todo.User{
				Name: "Test",
				Username: "test",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User){
				s.EXPECT().CreateUser(user).Return(1, nil)

			},
			expectedStatusCode: 200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name: "Empty Fields",
			inputBody: `{"username":"test","password":"qwerty"}`,
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User){},
			expectedStatusCode: 400,
			expectedRequestBody: `{"message":"Key: 'User.Name' Error:Field validation for 'Name' failed on the 'required' tag"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T){
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			service := &service.Service{Authorization: auth}
			handler := NewHandler(service)

			// Test Server
			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			// Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))

			// Perform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
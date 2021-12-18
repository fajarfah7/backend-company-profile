package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	// "time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"net/http"
	"net/http/httptest"
	"project_company_profile/models/admin"
	"project_company_profile/response"
)

// example to do request
// https://golang.hotexamples.com/examples/github.com.gin-gonic.gin/Engine/ServeHTTP/golang-engine-servehttp-method-examples.html

type MockUsecase struct {
	mock.Mock
}

func (mock *MockUsecase) CreateOrUpdate(textName, textType, newText string) (res response.Response) {
	args := mock.Called()

	if textName == "" {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusBadRequest
		res.Messages = []map[string]string{{"api": "Bad request, name is empty!"}}
		return res
	}
	if textType == "" {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusBadRequest
		res.Messages = []map[string]string{{"api": "Bad request, type is empty!"}}
		return res
	}
	if newText == "" {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusBadRequest
		res.Messages = []map[string]string{{"api": "Bad request, new text is empty!"}}
		return res
	}

	resp := args.Get(0)

	res = resp.(response.Response)

	return res
}

func (mock *MockUsecase) GetAll() (res response.Response) {
	args := mock.Called()
	resp := args.Get(0)
	data := resp.(response.Response)

	return data
}

func TestCreateOrUpdate(t *testing.T) {
	mockUcs := new(MockUsecase)

	textName := "test create or update"
	textType := "about_us_history"
	newText := "this is new text"

	resp := response.Response{
		ActionStatus: true,
		Data:         nil,
		Status:       http.StatusOK,
		Messages:     []map[string]string{{"api": "OK"}},
	}

	mockUcs.On("CreateOrUpdate", mock.Anything).Return(resp)

	aboutUsDelivery := NewAboutUsDelivery(mockUcs)

	type ReqParams struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Text string `json:"text"`
	}
	reqParams := ReqParams{
		Name: textName,
		Type: textType,
		Text: newText,
	}
	reqParamsJson, err := json.Marshal(reqParams)
	assert.NoError(t, err)

	server := gin.Default()
	server.POST("/about-us/save", aboutUsDelivery.CreateOrUpdate)

	payload := bytes.NewBufferString(string(reqParamsJson))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/about-us/save", payload)

	server.ServeHTTP(w, r)

	fmt.Println("payload ===>", payload)

	responseTemplate := new(response.Response)

	res := w.Result()
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, responseTemplate)
	assert.NoError(t, err)

	mockUcs.AssertExpectations(t)

	assert.Equal(t, true, responseTemplate.ActionStatus)

}

func TestGetAll(t *testing.T) {
	mockUcs := new(MockUsecase)

	timeNow := time.Now()
	aboutUs := []*admin.AboutUs{
		&admin.AboutUs{
			ID:        int64(1),
			ImageID:   int64(1),
			Name:      "History",
			Type:      "about_us_history",
			Text:      "Some text for history about us",
			CreatedAt: &timeNow,
			UpdatedAt: &timeNow,
		},
		&admin.AboutUs{
			ID:        int64(2),
			ImageID:   int64(1),
			Name:      "Vision",
			Type:      "about_us_vision",
			Text:      "Some text for vision about us",
			CreatedAt: &timeNow,
			UpdatedAt: &timeNow,
		},
		&admin.AboutUs{
			ID:        int64(2),
			ImageID:   int64(1),
			Name:      "Mission",
			Type:      "about_us_mission",
			Text:      "Some text for mission about us",
			CreatedAt: &timeNow,
			UpdatedAt: &timeNow,
		},
	}

	resp := response.Response{
		ActionStatus: true,
		Data:         aboutUs,
		Status:       http.StatusOK,
		Messages:     []map[string]string{{"api": "OK"}},
	}

	mockUcs.On("GetAll").Return(resp)

	aboutUsDelivery := NewAboutUsDelivery(mockUcs)

	server := gin.Default()
	server.POST("/about-us", aboutUsDelivery.GetAll)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/about-us", nil)
	w.Header().Set("Content-Type", "application/json")

	server.ServeHTTP(w, r)

	responseTemplate := new(response.Response)

	res := w.Result()
	body, _ := io.ReadAll(res.Body)
	err := json.Unmarshal(body, responseTemplate)
	assert.NoError(t, err)
	data := responseTemplate.Data.([]interface{})

	mockUcs.AssertExpectations(t)
	assert.Equal(t, len(aboutUs), len(data))
}

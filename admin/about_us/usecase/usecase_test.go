package usecase

import (
	"fmt"
	"project_company_profile/models/admin"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mock *MockRepository) CreateOrUpdate(textName, textType, newText string) error {
	args := mock.Called()
	fmt.Println("textName", textName)
	fmt.Println("textType", textType)
	fmt.Println("newText", newText)
	return args.Error(0)

}
func (mock *MockRepository) GetAll() ([]*admin.AboutUs, error) {
	args := mock.Called()
	resp := args.Get(0)
	return resp.([]*admin.AboutUs), args.Error(1)
}

func TestCreateOrUpdate(t *testing.T) {
	mockRepo := new(MockRepository)

	textName := "test create or update"
	textType := "test"
	newText := "this is new text"

	mockRepo.On("CreateOrUpdate", mock.Anything).Return(nil)

	aboutUsUsecase := NewAboutUsUsecase(mockRepo)

	createOrUpdate := aboutUsUsecase.CreateOrUpdate(textName, textType, newText)

	mockRepo.AssertExpectations(t)

	assert.Equal(t, 200, createOrUpdate.Status)
}

func TestGetAll(t *testing.T) {
	mockRepo := new(MockRepository)
	var timeNow time.Time
	timeNow = time.Now()
	aboutUs := &admin.AboutUs{
		ID:        1,
		ImageID:   1,
		Name:      "History",
		Type:      "about_us_history",
		Text:      "some text for about us history",
		CreatedAt: &timeNow,
		UpdatedAt: &timeNow,
	}
	mockRepo.On("GetAll").Return([]*admin.AboutUs{aboutUs}, nil)

	aboutUsUsecase := NewAboutUsUsecase(mockRepo)

	getAll := aboutUsUsecase.GetAll()

	response := getAll.Data.([]*admin.AboutUs)

	mockRepo.AssertExpectations(t)

	assert.Equal(t, int64(1), response[0].ID)
}

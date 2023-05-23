package mock

import (
	"github.com/acikkaynak/musahit-harita-backend/model/city"
	"github.com/acikkaynak/musahit-harita-backend/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

type MockedRepository struct {
	mock.Mock
}

func (m *MockedRepository) GetDistricts() []city.District {
	args := m.Called()
	return args.Get(0).([]city.District)
}

func TestGetFeeds(t *testing.T) {
	// Arrange
	mockRepository := new(MockedRepository)
	feeds := make([]city.District, 0)
	for i := 1; i <= 10; i++ {
		feeds = append(feeds, city.District{
			Id:   int64(i),
			Name: "District " + strconv.Itoa(i),
		})
	}

	mockRepository.On("GetDistricts").Return(feeds)

	repository.Districts = mockRepository.GetDistricts()

	// Act
	response, err := GetFeeds()

	// Assert
	require.NoError(t, err)
	require.Equal(t, 10, response.Count)
	require.Len(t, response.Results, 10)
	for _, result := range response.Results {
		require.Contains(t, []int{1, 2, 3, 4, 5}, result.VolunteerData)
	}
	mockRepository.AssertExpectations(t)
}

func TestGetFeedsWithEmptyDistricts(t *testing.T) {
	// Arrange
	mockRepository := new(MockedRepository)
	feeds := make([]city.District, 0)

	mockRepository.On("GetDistricts").Return(feeds)

	repository.Districts = mockRepository.GetDistricts()

	// Act
	response, err := GetFeeds()

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, response.Count)
	require.Len(t, response.Results, 0)
	mockRepository.AssertExpectations(t)
}

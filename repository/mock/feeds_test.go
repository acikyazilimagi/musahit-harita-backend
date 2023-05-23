package mock

import (
	"github.com/acikkaynak/musahit-harita-backend/model"
	"github.com/acikkaynak/musahit-harita-backend/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

type MockedRepository struct {
	mock.Mock
}

func (m *MockedRepository) GetNeighborhoods() map[int]model.Neighborhood {
	args := m.Called()
	return args.Get(0).(map[int]model.Neighborhood)
}

func TestGetFeeds(t *testing.T) {
	// Arrange
	mockRepository := new(MockedRepository)
	feeds := make(map[int]model.Neighborhood)
	for i := 1; i <= 10; i++ {
		feeds[i] = model.Neighborhood{
			Id:   i,
			Name: "Neighborhood " + strconv.Itoa(i),
		}
	}

	mockRepository.On("GetNeighborhoods").Return(feeds)

	repository.NeighborhoodIdToMap = mockRepository.GetNeighborhoods()

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
	feeds := make(map[int]model.Neighborhood)

	mockRepository.On("GetNeighborhoods").Return(feeds)

	repository.NeighborhoodIdToMap = mockRepository.GetNeighborhoods()

	// Act
	response, err := GetFeeds()

	// Assert
	require.NoError(t, err)
	require.Equal(t, 0, response.Count)
	require.Len(t, response.Results, 0)
	mockRepository.AssertExpectations(t)
}

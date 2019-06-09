package controller_test

import (
	"encoding/json"
	"github.com/fernetbalboa/arqweb/config"
	"github.com/fernetbalboa/arqweb/controller"
	"github.com/fernetbalboa/arqweb/mock"
	"github.com/fernetbalboa/arqweb/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	config.Scope = config.TestScope
}

func TestAddPOI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	addPOIEndpoint := "/poi"

	t.Run("Request is not GeoJson", func(t *testing.T) {
		//Given
		POIController := &controller.POIController{}
		badRequestBody := `{"invalidField": "not a GeoJson"}`

		r := config.ConfiguredRouter()
		r.POST(addPOIEndpoint, POIController.AddPOI)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", addPOIEndpoint, strings.NewReader(badRequestBody))

		//When
		r.ServeHTTP(w, req)

		//Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, `{"message":"There was a problem with the request","error":"bad_request","status":400,"cause":["Error parsing POI. It should be a GeoJson feature"]}`,
			w.Body.String())

	})

	t.Run("Add POI successfully", func(t *testing.T) {
		//Given
		storageMock := mock.NewMockPOIStorage(ctrl)
		POIController := &controller.POIController{POIStorage: storageMock}

		r := config.ConfiguredRouter()
		r.POST(addPOIEndpoint, POIController.AddPOI)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", addPOIEndpoint, strings.NewReader(test.DefaultGeoJsonFeatureString()))

		savedPOI := test.DefaultPOI()
		storageMock.EXPECT().SaveFeature(test.DefaultGeoJsonFeature()).Return(savedPOI, nil)

		//When
		r.ServeHTTP(w, req)

		//Then
		assert.Equal(t, http.StatusCreated, w.Code)
		respBody, _ := json.Marshal(savedPOI)

		assert.Equal(t, respBody, w.Body.Bytes())

	})
}

/*
func TestSearchCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	searchPOIEndpoint := "/category/search"

	t.Run("Search POI successfully", func(t *testing.T) {
		//Given
		storageMock := mock.NewMockPOIStorage(ctrl)
		POIController := &controller.POIController{POIStorage: storageMock}

		category := "food"
		maxLat := 10.63462
		minLat := 5.45345
		maxLong := 30.3523523
		minLong := 24.62311424

		queryParams := url.Values{}
		queryParams.Add("category", category)
		queryParams.Add("maxLat", strconv.FormatFloat(maxLat, 'f', 10, 64))
		queryParams.Add("minLat", strconv.FormatFloat(minLat, 'f', 10, 64))
		queryParams.Add("maxLong", strconv.FormatFloat(maxLong, 'f', 10, 64))
		queryParams.Add("minLong", strconv.FormatFloat(minLong, 'f', 10, 64))
		queryParams.Add("bound", strconv.FormatBool(true))
		queryParams.Add("hidden", strconv.FormatBool(false))

		r := config.ConfiguredRouter()
		r.GET(searchPOIEndpoint, POIController.SearchPOI)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", searchPOIEndpoint, nil)
		req.URL.RawQuery = queryParams.Encode()

		savedPOI1 := test.DefaultPOI()
		savedPOI2 := test.DefaultPOI()
		savedPOI2.Id = test.NewDocumentId()
		savedPOI3 := test.DefaultPOI()
		savedPOI3.Hidden = true
		savedPOI3.Id = test.NewDocumentId()

		foundPOIs := []*domain.PointOfInterest{savedPOI1, savedPOI2}

		filters := &domain.POIFilter{
			Category: category,
			MaxLong:  maxLong,
			MaxLat:   maxLat,
			MinLat:   minLat,
			MinLong:  minLong,
			Bound:    true,
			Limit:    controller.DefaultSearchLimit,
			Hidden:   false,
		}

		storageMock.EXPECT().SearchPOI(filters).Return(foundPOIs, nil)

		//When
		r.ServeHTTP(w, req)

		//Then
		assert.Equal(t, http.StatusOK, w.Code)
		respBody, _ := json.Marshal(foundPOIs)

		assert.Equal(t, respBody, w.Body.Bytes())

	})
}
*/

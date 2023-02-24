package v1_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/malkev1ch/apod/gen/v1"
	"github.com/malkev1ch/apod/internal/constant"
	v1 "github.com/malkev1ch/apod/internal/handler/v1"
	"github.com/malkev1ch/apod/internal/handler/v1/mock"
	"github.com/malkev1ch/apod/internal/model"
	"github.com/malkev1ch/apod/pkg/pointer"
	"github.com/stretchr/testify/require"
)

const dayDuration = time.Hour * 24

func TestPicture_GetPicture(t *testing.T) {
	truncatedNow := time.Now().UTC().Truncate(dayDuration)
	expected := fmt.Sprintf(`{"date":"%s","title":"Title","url":"LocalURL"}`, truncatedNow.Format(constant.DateFormatISO))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctl := gomock.NewController(t)
	svc := mock.NewMockPictureService(ctl)
	svc.EXPECT().GetByDate(gomock.Any(), gomock.Any()).
		Return(
			model.Picture{
				Date:        truncatedNow,
				Title:       "Title",
				URL:         "URL",
				HDURL:       pointer.New("HDURL"),
				ThumbURL:    nil,
				LocalURL:    "LocalURL",
				MediaType:   "MediaType",
				Copyright:   pointer.New("Copyright"),
				Explanation: "Explanation",
			},
			nil)

	handler := v1.NewPicture(svc)
	err := handler.GetPicture(c, gen.GetPictureParams{Date: &openapi_types.Date{Time: truncatedNow}})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, expected, strings.TrimSuffix(rec.Body.String(), "\n"))
}

func TestPicture_GetAllPictures(t *testing.T) {
	truncatedNow := time.Now().UTC().Truncate(dayDuration)
	expected := fmt.Sprintf(
		`[{"date":"%s","title":"Title 1","url":"LocalURL 1"},{"date":"%s","title":"Title 2","url":"LocalURL 2"}]`,
		truncatedNow.Format(constant.DateFormatISO),
		truncatedNow.Add(dayDuration).Format(constant.DateFormatISO))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	ctl := gomock.NewController(t)
	svc := mock.NewMockPictureService(ctl)
	svc.EXPECT().GetAll(gomock.Any()).
		Return(
			[]model.Picture{
				{
					Date:        truncatedNow,
					Title:       "Title 1",
					URL:         "URL 1",
					HDURL:       pointer.New("HDURL 1"),
					ThumbURL:    nil,
					LocalURL:    "LocalURL 1",
					MediaType:   "MediaType 1",
					Copyright:   pointer.New("Copyright 1"),
					Explanation: "Explanation 1",
				},
				{
					Date:        truncatedNow.Add(dayDuration),
					Title:       "Title 2",
					URL:         "URL 2",
					HDURL:       pointer.New("HDURL 2"),
					ThumbURL:    nil,
					LocalURL:    "LocalURL 2",
					MediaType:   "MediaType 2",
					Copyright:   pointer.New("Copyright 2"),
					Explanation: "Explanation 2",
				},
			},
			nil)

	handler := v1.NewPicture(svc)
	err := handler.GetAllPictures(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, expected, strings.TrimSuffix(rec.Body.String(), "\n"))
}

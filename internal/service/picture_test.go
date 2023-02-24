package service_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/malkev1ch/apod/internal/model"
	"github.com/malkev1ch/apod/internal/service"
	"github.com/malkev1ch/apod/internal/service/mock"
	berror "github.com/malkev1ch/apod/pkg/errors"
	"github.com/malkev1ch/apod/pkg/pointer"
	"github.com/stretchr/testify/require"
)

const dayDuration = time.Hour * 24

func TestPicture_GetByDate(t *testing.T) {
	ctx := context.Background()
	truncatedNow := time.Now().UTC().Truncate(dayDuration)
	bucket := "bucket"
	tests := []struct {
		name        string
		expected    model.Picture
		pictureRepo func(*gomock.Controller) service.PictureRepository
		nasaAPIRepo func(*gomock.Controller) service.NasaAPIRepository
		fileRepo    func(*gomock.Controller) service.FileRepository
		wantErr     bool
	}{
		{
			name: "Already exists",
			expected: model.Picture{
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
			pictureRepo: func(ctl *gomock.Controller) service.PictureRepository {
				m := mock.NewMockPictureRepository(ctl)
				m.EXPECT().GetByDate(ctx, truncatedNow).Return(
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

				return m
			},
			nasaAPIRepo: func(ctl *gomock.Controller) service.NasaAPIRepository {
				return mock.NewMockNasaAPIRepository(ctl)
			},
			fileRepo: func(ctl *gomock.Controller) service.FileRepository {
				return mock.NewMockFileRepository(ctl)
			},
			wantErr: false,
		},
		{
			name: "Doesn't exist",
			expected: model.Picture{
				Date:        truncatedNow,
				Title:       "Title",
				URL:         "URL",
				HDURL:       pointer.New("HDURL"),
				ThumbURL:    nil,
				LocalURL:    "LocalURL",
				MediaType:   "image",
				Copyright:   pointer.New("Copyright"),
				Explanation: "Explanation",
			},
			pictureRepo: func(ctl *gomock.Controller) service.PictureRepository {
				m := mock.NewMockPictureRepository(ctl)
				m.EXPECT().GetByDate(
					ctx,
					truncatedNow,
				).Return(model.Picture{}, berror.New("", "", berror.PictureNotFoundErrCode))
				m.EXPECT().Create(
					ctx,
					&model.Picture{
						Date:        truncatedNow,
						Title:       "Title",
						URL:         "URL",
						HDURL:       pointer.New("HDURL"),
						ThumbURL:    nil,
						LocalURL:    "LocalURL",
						MediaType:   "image",
						Copyright:   pointer.New("Copyright"),
						Explanation: "Explanation",
					},
				).Return(nil)
				return m
			},
			nasaAPIRepo: func(ctl *gomock.Controller) service.NasaAPIRepository {
				m := mock.NewMockNasaAPIRepository(ctl)
				m.EXPECT().GetAPODByDate(ctx, truncatedNow).
					Return(
						model.Picture{
							Date:        truncatedNow,
							Title:       "Title",
							URL:         "URL",
							HDURL:       pointer.New("HDURL"),
							ThumbURL:    nil,
							LocalURL:    "",
							MediaType:   "image",
							Copyright:   pointer.New("Copyright"),
							Explanation: "Explanation",
						},
						nil)
				m.EXPECT().LoadFile(ctx, "HDURL").
					Return(
						model.File{
							Size:        10,
							Name:        "Name",
							ContentType: "ContentType",
							File:        bytes.NewReader([]byte{}),
						},
						nil)
				return m
			},
			fileRepo: func(ctl *gomock.Controller) service.FileRepository {
				m := mock.NewMockFileRepository(ctl)
				m.EXPECT().PutObject(
					ctx,
					bucket,
					&model.File{
						Size:        10,
						Name:        "Name",
						ContentType: "ContentType",
						File:        bytes.NewReader([]byte{}),
					},
				).Return("LocalURL", nil)
				return m
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			svc := service.NewPicture(tt.nasaAPIRepo(ctl), tt.pictureRepo(ctl), tt.fileRepo(ctl), bucket)
			actual, err := svc.GetByDate(ctx, truncatedNow)
			if tt.wantErr {
				require.Equal(t, err != nil, tt.wantErr)
				return
			}
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestPicture_GetAll(t *testing.T) {
	ctx := context.Background()
	truncatedNow := time.Now().UTC().Truncate(dayDuration)
	bucket := "bucket"

	expected := []model.Picture{
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
			HDURL:       nil,
			ThumbURL:    pointer.New("ThumbURL 2"),
			LocalURL:    "LocalURL 2",
			MediaType:   "MediaType 2",
			Copyright:   pointer.New("Copyright 2"),
			Explanation: "Explanation 2",
		},
	}

	ctl := gomock.NewController(t)
	nasaAPIRepo := mock.NewMockNasaAPIRepository(ctl)
	fileRepo := mock.NewMockFileRepository(ctl)
	pictureRepo := mock.NewMockPictureRepository(ctl)
	pictureRepo.EXPECT().GetAll(ctx).Return(expected, nil)

	svc := service.NewPicture(nasaAPIRepo, pictureRepo, fileRepo, bucket)

	actual, err := svc.GetAll(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t, expected, actual)
}

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/malkev1ch/apod/internal/model"
	"github.com/malkev1ch/apod/internal/repository"
	"github.com/malkev1ch/apod/pkg/pointer"
	"github.com/stretchr/testify/require"
)

func TestPicture_Create(t *testing.T) {
	ctx := context.Background()

	input := model.Picture{
		Date:        time.Now().UTC(),
		Title:       "Title",
		URL:         "URL",
		HDURL:       pointer.New("HDURL"),
		ThumbURL:    nil,
		LocalURL:    "LocalURL",
		MediaType:   "MediaType",
		Copyright:   pointer.New("Copyright"),
		Explanation: "Explanation",
	}

	repo := repository.NewPicture(pool)

	err := repo.Create(ctx, &input)
	require.NoError(t, err)

	require.NoError(t, clearDatabase(ctx))
}

func TestPicture_GetByDate(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	truncatedNow := now.Truncate(dayDuration)
	tests := []struct {
		name     string
		date     time.Time
		input    model.Picture
		expected model.Picture
		wantErr  bool
	}{
		{
			name: "Success",
			date: now,
			input: model.Picture{
				Date:        now,
				Title:       "Title",
				URL:         "URL",
				HDURL:       pointer.New("HDURL"),
				ThumbURL:    nil,
				LocalURL:    "LocalURL",
				MediaType:   "MediaType",
				Copyright:   pointer.New("Copyright"),
				Explanation: "Explanation",
			},
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
			wantErr: false,
		},
		{
			name: "Not found",
			date: truncatedNow.Add(dayDuration),
			input: model.Picture{
				Date:        now,
				Title:       "Title",
				URL:         "URL",
				HDURL:       pointer.New("HDURL"),
				ThumbURL:    nil,
				LocalURL:    "LocalURL",
				MediaType:   "MediaType",
				Copyright:   pointer.New("Copyright"),
				Explanation: "Explanation",
			},
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
			wantErr: true,
		},
	}

	repo := repository.NewPicture(pool)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, &tt.input)
			require.NoError(t, err)

			actual, err := repo.GetByDate(ctx, tt.date)

			if tt.wantErr {
				require.Equal(t, err != nil, tt.wantErr)
				require.NoError(t, clearDatabase(ctx))
				return
			}

			require.Equal(t, tt.expected, actual)

			require.NoError(t, clearDatabase(ctx))
		})
	}
}

func TestPicture_GetAll(t *testing.T) {
	ctx := context.Background()
	truncatedNow := time.Now().UTC().Truncate(dayDuration)

	testPictures := []model.Picture{
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

	repo := repository.NewPicture(pool)

	for i := range testPictures {
		err := repo.Create(ctx, &testPictures[i])
		require.NoError(t, err)
	}

	actual, err := repo.GetAll(ctx)
	require.NoError(t, err)

	require.ElementsMatch(t, testPictures, actual)

	require.NoError(t, clearDatabase(ctx))
}

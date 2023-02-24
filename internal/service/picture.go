package service

import (
	"context"
	"time"

	"github.com/malkev1ch/apod/internal/model"
	berror "github.com/malkev1ch/apod/pkg/errors"
	"github.com/malkev1ch/apod/pkg/logger"
)

//go:generate mockgen -source=picture.go -destination=./mock/picture_mock.go -package=mock

type PictureRepository interface {
	Create(context.Context, *model.Picture) error
	GetByDate(context.Context, time.Time) (model.Picture, error)
	GetAll(context.Context) ([]model.Picture, error)
}

type NasaAPIRepository interface {
	GetAPODByDate(ctx context.Context, date time.Time) (p model.Picture, err error)
	LoadFile(ctx context.Context, rawURL string) (f model.File, err error)
}

type FileRepository interface {
	PutObject(ctx context.Context, bucket string, input *model.File) (string, error)
}

// Picture represents picture behavior.
type Picture struct {
	nasaAPI     NasaAPIRepository
	pictureRepo PictureRepository
	FileRepo    FileRepository
	bucket      string
}

// NewPicture creates a new Picture.
func NewPicture(
	nasaAPI NasaAPIRepository,
	pictureRepo PictureRepository,
	fileRepo FileRepository,
	bucket string,
) *Picture {
	return &Picture{
		nasaAPI:     nasaAPI,
		pictureRepo: pictureRepo,
		FileRepo:    fileRepo,
		bucket:      bucket,
	}
}

func (p *Picture) GetByDate(ctx context.Context, date time.Time) (model.Picture, error) {
	pic, err := p.pictureRepo.GetByDate(ctx, date)
	if err != nil {
		if berror.Is(err, berror.PictureNotFoundErrCode) {
			nPic, inErr := p.nasaAPI.GetAPODByDate(ctx, date)
			if inErr != nil {
				return model.Picture{}, berror.ExtendPath("GetByDate-p.nasaAPI.GetAPODByDate", inErr)
			}

			urlToLoad := nPic.ChooseURL()
			if urlToLoad == "" {
				return model.Picture{},
					// business err for frontend.
					berror.New(
						"nPic.ChooseUrl",
						"Invalid picture format in NASA.",
						berror.NasaPictureInvalidFormatErrCode,
					)
			}

			f, inErr := p.nasaAPI.LoadFile(ctx, urlToLoad)
			if inErr != nil {
				return model.Picture{}, berror.ExtendPath("GetByDate-p.nasaAPI.LoadFile", inErr)
			}

			localURL, inErr := p.FileRepo.PutObject(ctx, p.bucket, &f)
			if inErr != nil {
				return model.Picture{}, berror.ExtendPath("GetByDate-p.FileRepo.PutObject", inErr)
			}

			logger.LoggerFromContext(ctx).Debugw(
				"successfully loaded file",
				"urlToLoad", urlToLoad,
				"newURL", localURL,
			)

			nPic.LocalURL = localURL

			inErr = p.pictureRepo.Create(ctx, &nPic)
			if inErr != nil {
				return model.Picture{}, berror.ExtendPath("GetByDate-p.pictureRepo.Create", inErr)
			}

			return nPic, nil
		}

		return model.Picture{}, berror.ExtendPath("GetByDate-p.pictureRepo.GetByDate", err)
	}

	return pic, nil
}

func (p *Picture) GetAll(ctx context.Context) ([]model.Picture, error) {
	pics, err := p.pictureRepo.GetAll(ctx)
	if err != nil {
		return nil, berror.ExtendPath("GetAll-p.pictureRepo.GetAll", err)
	}

	return pics, nil
}

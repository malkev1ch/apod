package v1

import (
	"context"
	"net/http"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/labstack/echo/v4"
	"github.com/malkev1ch/apod/gen/v1"
	"github.com/malkev1ch/apod/internal/model"
	berror "github.com/malkev1ch/apod/pkg/errors"
	"github.com/malkev1ch/apod/pkg/logger"
)

//go:generate mockgen -source=picture.go -destination=./mock/picture_mock.go -package=mock

type PictureService interface {
	GetByDate(ctx context.Context, date time.Time) (model.Picture, error)
	GetAll(ctx context.Context) ([]model.Picture, error)
}

// Picture represents picture handling layer.
type Picture struct {
	svc PictureService
}

// NewPicture creates a new Picture.
func NewPicture(svc PictureService) *Picture {
	return &Picture{svc: svc}
}

// GetPicture handles GET {baseUrl}/v1/picture method.
func (p *Picture) GetPicture(ctx echo.Context, params gen.GetPictureParams) error {
	t := time.Now().UTC()
	if params.Date != nil {
		t = params.Date.Time
	}

	pic, err := p.svc.GetByDate(logger.ContextEchoPropagateLogger(ctx), t)
	if err != nil {
		return berror.ExtendPath("GetPicture-p.svc.GetByDate", err)
	}

	resp := p.toGenPicture(&pic)

	return ctx.JSON(http.StatusOK, &resp)
}

// GetAllPictures handles GET {baseUrl}/v1/picture/all method.
func (p *Picture) GetAllPictures(ctx echo.Context) error {
	pics, err := p.svc.GetAll(logger.ContextEchoPropagateLogger(ctx))
	if err != nil {
		return berror.ExtendPath("GetAllPicture-p.svc.GetAll", err)
	}

	resp := make(gen.PictureArray, 0, len(pics))
	for i := range pics {
		resp = append(resp, p.toGenPicture(&pics[i]))
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (p *Picture) toGenPicture(in *model.Picture) gen.Picture {
	return gen.Picture{
		Date:  openapi_types.Date{Time: in.Date},
		Title: in.Title,
		Url:   in.LocalURL,
	}
}

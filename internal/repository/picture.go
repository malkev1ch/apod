package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malkev1ch/apod/internal/model"
	berror "github.com/malkev1ch/apod/pkg/errors"
)

type Picture struct {
	pool *pgxpool.Pool
}

// NewPicture creates a new Picture.
func NewPicture(pool *pgxpool.Pool) *Picture {
	return &Picture{pool: pool}
}

// Create creates a new picture in postgres.
func (p *Picture) Create(ctx context.Context, pic *model.Picture) error {
	const q = `
			INSERT INTO pictures
				(id,"date",title,url,local_url,hd_url,thumbnail_url,media_type,copyright,explanation,created_at)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := p.pool.Exec(
		ctx,
		q,
		uuid.New(),
		pic.Date,
		pic.Title,
		pic.URL,
		pic.LocalURL,
		pic.HDURL,
		pic.ThumbURL,
		pic.MediaType,
		pic.Copyright,
		pic.Explanation,
		time.Now().UTC(),
	)
	if err != nil {
		return berror.NewInternal("Create-p.pool.Exec", err)
	}

	return nil
}

// GetByDate returns a picture by date from postgres.
func (p *Picture) GetByDate(ctx context.Context, date time.Time) (pic model.Picture, err error) {
	const q = `
			SELECT "date",
			       title,
			       url,
				   local_url,
			       hd_url,
			       thumbnail_url,
			       media_type,
			       copyright,
			       explanation
			FROM pictures
			WHERE date = $1`

	row := p.pool.QueryRow(ctx, q, date)
	pic, err = p.scanRow(row)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return model.Picture{}, berror.New(
				"GetByDate-p.ScanRow",
				"Picture not found in local repository",
				berror.PictureNotFoundErrCode,
			)
		}
		return model.Picture{}, berror.NewInternal("GetByDate-p.scanRow", err)
	}

	return pic, nil
}

// GetAll returns all pictures postgres.
func (p *Picture) GetAll(ctx context.Context) (pics []model.Picture, err error) {
	const q = `
			SELECT "date",
			       title,
			       url,
				   local_url,
			       hd_url,
			       thumbnail_url,
			       media_type,
			       copyright,
			       explanation
			FROM pictures`

	rows, err := p.pool.Query(ctx, q)
	if err != nil {
		return nil, berror.NewInternal("GetAll-p.pool.Query", err)
	}
	defer rows.Close()

	for rows.Next() {
		pic, inErr := p.scanRow(rows)
		if inErr != nil {
			return nil, berror.NewInternal("GetAll-p.scanRow", err)
		}

		pics = append(pics, pic)
	}

	return pics, nil
}

func (p *Picture) scanRow(row pgx.Row) (pic model.Picture, err error) {
	err = row.Scan(
		&pic.Date,
		&pic.Title,
		&pic.URL,
		&pic.LocalURL,
		&pic.HDURL,
		&pic.ThumbURL,
		&pic.MediaType,
		&pic.Copyright,
		&pic.Explanation,
	)
	if err != nil {
		return model.Picture{}, err
	}

	return pic, err
}

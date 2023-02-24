package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/malkev1ch/apod/internal/constant"
	"github.com/malkev1ch/apod/internal/model"
)

const (
	nasaAPIDateParam   = "date"
	nasaAPIThumbsParam = "thumbs"
	nasaAPIHDParam     = "hd"
	nasaAPIApiKeyParam = "api_key"
	// ContentTypeImageXYZ is generic content-type for images
	ContentTypeImageXYZ = "image/xyz"
)

// NasaAPI = http client to Nasa API.
type NasaAPI struct {
	httpClient *http.Client
	apiKey     string
	addr       string
}

// NewNasaAPI creates a new NasaAPI.
func NewNasaAPI(httpClient *http.Client, apiKey, addr string) *NasaAPI {
	return &NasaAPI{
		httpClient: httpClient,
		apiKey:     apiKey,
		addr:       addr,
	}
}

type getAPODByDateResponse struct {
	Date        string  `json:"date"`
	Title       string  `json:"title"`
	URL         string  `json:"url"`
	HDURL       *string `json:"hd_url"`
	ThumbURL    *string `json:"thumbnail_url"`
	MediaType   string  `json:"media_type"`
	Copyright   *string `json:"copyright"`
	Explanation string  `json:"explanation"`
}

func (n *NasaAPI) GetAPODByDate(ctx context.Context, date time.Time) (p model.Picture, err error) {
	u := n.buildGetAPODByDateURL(date)
	req, err := http.NewRequestWithContext(ctx, "GET", u, http.NoBody)
	if err != nil {
		return model.Picture{}, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return model.Picture{}, err
	}
	defer func() {
		if bodyErr := resp.Body.Close(); bodyErr != nil {
			err = bodyErr
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Picture{}, fmt.Errorf("io.ReadAll: %w", err)
	}

	var r getAPODByDateResponse

	err = json.Unmarshal(body, &r)
	if err != nil {
		return model.Picture{}, fmt.Errorf("json.Unmarshal: %w", err)
	}

	t, err := time.Parse(constant.DateFormatISO, r.Date)
	if err != nil {
		return model.Picture{}, fmt.Errorf("time.Parse: %w", err)
	}

	p = model.Picture{
		Date:        t,
		Title:       r.Title,
		URL:         r.URL,
		HDURL:       r.HDURL,
		ThumbURL:    r.ThumbURL,
		LocalURL:    "",
		MediaType:   r.MediaType,
		Copyright:   r.Copyright,
		Explanation: r.Explanation,
	}

	return p, nil
}

func (n *NasaAPI) buildGetAPODByDateURL(date time.Time) string {
	var u url.URL
	u.Scheme = "https"
	u.Host = n.addr
	u.Path = "/planetary/apod"
	q := u.Query()
	q.Set(nasaAPIApiKeyParam, n.apiKey)
	q.Set(nasaAPIDateParam, date.Format(constant.DateFormatISO))
	q.Set(nasaAPIThumbsParam, "true")
	q.Set(nasaAPIHDParam, "true")
	u.RawQuery = q.Encode()
	return u.String()
}

// LoadFile runs file download by url.
func (n *NasaAPI) LoadFile(ctx context.Context, rawURL string) (f model.File, err error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return model.File{}, fmt.Errorf("failed to parse provided image url %s: %w", rawURL, err)
	}

	imageName := fmt.Sprintf("%s_%s", uuid.New().String(), path.Base(parsedURL.Path))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, http.NoBody)
	if err != nil {
		return model.File{}, fmt.Errorf("http.NewRequestWithContext: %w", err)
	}

	res, err := n.httpClient.Do(req)
	if err != nil {
		return model.File{}, fmt.Errorf("n.httpClient.Do: %w", err)
	}
	defer func() {
		if bodyErr := res.Body.Close(); bodyErr != nil {
			err = bodyErr
		}
	}()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return model.File{}, fmt.Errorf("io.ReadAll: %w", err)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType == "" {
		contentType = ContentTypeImageXYZ
	}

	return model.File{
		Size:        res.ContentLength,
		Name:        imageName,
		ContentType: contentType,
		File:        bytes.NewReader(content),
	}, nil
}

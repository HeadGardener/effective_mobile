package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/HeadGardener/effective_mobile/internal/config"
)

type Client struct {
	cl                 *http.Client
	ageBaseURL         string
	genderBaseURL      string
	nationalityBaseURL string
}

func NewClient(conf config.HTTPClientConfig) *Client {
	return &Client{
		cl:                 http.DefaultClient,
		ageBaseURL:         conf.AgeBaseURL,
		genderBaseURL:      conf.GenderBaseURL,
		nationalityBaseURL: conf.NationalityBaseURL,
	}
}

func (c *Client) GetAge(ctx context.Context, name string) (int8, error) {
	resp, err := c.sendGetRequest(ctx, c.ageBaseURL+"?name="+name)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return 0, err
	}

	type ageResp struct {
		Age int8 `json:"age"`
	}

	var age ageResp
	if err = json.NewDecoder(resp.Body).Decode(&age); err != nil {
		return 0, err
	}

	return age.Age, nil
}

func (c *Client) GetGender(ctx context.Context, name string) (string, error) {
	resp, err := c.sendGetRequest(ctx, c.genderBaseURL+"?name="+name)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return "", err
	}

	type genderResp struct {
		Gender string `json:"gender"`
	}

	var gender genderResp
	if err = json.NewDecoder(resp.Body).Decode(&gender); err != nil {
		return "", err
	}

	return gender.Gender, nil
}

func (c *Client) GetNationality(ctx context.Context, name string) (string, error) {
	resp, err := c.sendGetRequest(ctx, c.nationalityBaseURL+"?name="+name)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return "", err
	}

	type Country struct {
		CountryID string `json:"country_id"`
	}
	type nationalityResp struct {
		Country []Country `json:"country"`
	}

	var nationality nationalityResp
	if err = json.NewDecoder(resp.Body).Decode(&nationality); err != nil {
		return "", err
	}

	return nationality.Country[0].CountryID, nil
}

func (c *Client) sendGetRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewBuffer([]byte("")))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := c.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}

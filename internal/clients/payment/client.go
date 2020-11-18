package payment

import (
	"encoding/json"
	"errors"
	"github.com/Solar-2020/SolAr_Backend_2020/internal/models"
	"github.com/valyala/fasthttp"
)

type client struct {
	host   string
	secret string
}

func NewClient(host, secret string) *client {
	return &client{
		host:   host,
		secret: secret,
	}
}

func (c *client) Create(createRequest models.CreateRequest) (createdPayments []models.Payment, err error) {
	if len(createRequest.Payments) == 0 {
		return
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().SetScheme("http")
	req.URI().SetHost(c.host)
	req.URI().SetPath("api/internal/payment/payment")

	req.Header.Set("Authorization", c.secret)
	req.Header.SetMethod(fasthttp.MethodPost)

	body, err := json.Marshal(createRequest)
	if err != nil {
		return
	}

	req.SetBody(body)

	err = fasthttp.Do(req, resp)
	if err != nil {
		return
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		err = json.Unmarshal(resp.Body(), &createdPayments)
		if err != nil {
			return
		}
		return
	case fasthttp.StatusBadRequest:
		var httpErr httpError
		err = json.Unmarshal(resp.Body(), &httpErr)
		if err != nil {
			return
		}
		return createdPayments, errors.New(httpErr.Error)

	default:
		return createdPayments, errors.New("Unexpected Server Error")
	}
}

func (c *client) GetByPostIDs(postIDs []int) (payments []models.Payment, err error) {
	payments = make([]models.Payment, 0)
	if len(postIDs) == 0 {
		return
	}

	ids := struct {
		PostIDs []int `json:"postIDs"`
	}{PostIDs: postIDs}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.URI().SetScheme("http")
	req.URI().SetHost(c.host)
	req.URI().SetPath("api/internal/payment/by-post-ids")

	req.Header.Set("Authorization", c.secret)
	req.Header.SetMethod(fasthttp.MethodPost)

	body, err := json.Marshal(ids)
	if err != nil {
		return
	}

	req.SetBody(body)

	err = fasthttp.Do(req, resp)
	if err != nil {
		return
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		err = json.Unmarshal(resp.Body(), &payments)
		if err != nil {
			return
		}
		return
	case fasthttp.StatusBadRequest:
		var httpErr httpError
		err = json.Unmarshal(resp.Body(), &httpErr)
		if err != nil {
			return
		}
		return payments, errors.New(httpErr.Error)

	default:
		return payments, errors.New("Unexpected Server Error")
	}
}

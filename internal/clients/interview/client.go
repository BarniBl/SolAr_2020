package interview

import (
	"bytes"
	"encoding/json"
	"errors"
	interviewModels "github.com/Solar-2020/Interview-Backend/pkg/models"
	"github.com/Solar-2020/SolAr_Backend_2020/internal/models"
	"net/http"
)

type Client interface {
	InsertInterviews(interviews []models.Interview, postID int) (err error)
	SelectInterviewsResults(postIDs []int, userID int) (interviews []interviewModels.InterviewResult, err error)
}

type client struct {
	host   string
	secret string
}

func NewClient(host string, secret string) Client {
	return &client{host: host, secret: secret}
}

type httpError struct {
	Error string `json:"error"`
}

type CreateRequest struct {
	Interviews []models.Interview `json:"interviews" validate:"required"`
	PostID     int                `json:"postID" validate:"required"`
}

func (c *client) InsertInterviews(interviews []models.Interview, postID int) (err error) {
	if len(interviews) == 0 {
		return
	}

	createRequest := CreateRequest{Interviews: interviews, PostID: postID}
	body, err := json.Marshal(createRequest)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, c.host+"/api/interview/create", bytes.NewReader(body))
	if err != nil {
		return
	}

	req.Header.Set("Authorization", c.secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		var httpErr httpError
		err = json.NewDecoder(resp.Body).Decode(&httpErr)
		if err != nil {
			return
		}
		return errors.New(httpErr.Error)
	default:
		return errors.New("Unexpected Server Error")
	}
}

type GetUniversalRequest struct {
	PostIDs         []int `json:"postIDs" validate:"required"`
	UserID          int   `json:"userID" validate:"required"`
	NotPassedResult bool  `json:"notPassedResult"`
}

type GetUniversalResponse struct {
	Interviews []interviewModels.InterviewResult `json:"interviews"`
}

func (c *client) SelectInterviewsResults(postIDs []int, userID int) (interviews []interviewModels.InterviewResult, err error) {
	if len(postIDs) == 0 {
		return
	}

	getUniversalRequest := GetUniversalRequest{PostIDs: postIDs, UserID: userID}
	body, err := json.Marshal(getUniversalRequest)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, c.host+"/api/interview/list", bytes.NewReader(body))
	if err != nil {
		return
	}

	req.Header.Set("Authorization", c.secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var response GetUniversalResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		return response.Interviews, err
	case http.StatusBadRequest:
		var httpErr httpError
		err = json.NewDecoder(resp.Body).Decode(&httpErr)
		if err != nil {
			return
		}
		return interviews, errors.New(httpErr.Error)
	default:
		return interviews, errors.New("Unexpected Server Error")
	}
}

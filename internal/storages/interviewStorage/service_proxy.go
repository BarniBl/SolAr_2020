package interviewStorage

import (
	service "github.com/Solar-2020/GoUtils/http"
	interviewApi "github.com/Solar-2020/Interview-Backend/pkg/api"
	interviewModels "github.com/Solar-2020/Interview-Backend/pkg/models"
	"github.com/Solar-2020/SolAr_Backend_2020/internal/models"
)

type ServiceStorage struct {
	service.Service
	address string
}

func (s *ServiceStorage) Address() string {
	return s.address
}

func NewStorageProxy(serviceAddress string) Storage {
	return &ServiceStorage{
		address: serviceAddress,
	}
}

func (s *ServiceStorage) InsertInterviews(interviews []models.Interview, postID int) (err error) {
	if len(interviews) == 0 {
		return
	}

	endpoint := service.ServiceEndpoint{
		Service:     s,
		Endpoint:    "/api/interview/create",
		Method:      "POST",
		ContentType: "application/json",
	}

	message := interviewApi.CreateRequest{
		Interviews: ToApiInterviews(interviews),
		PostID:     postID,
	}

	resp := interviewApi.CreateResponse{}
	err = endpoint.Send(message, &resp)
	if err != nil {
		return
	}
	interviews = FromApiInterviews(resp.Interviews)
	return err

}

func (s *ServiceStorage) SelectInterviews(postIDs []int) (interviews []models.Interview, err error) {
	endpoint := service.ServiceEndpoint{
		Service:     s,
		Endpoint:    "/api/interview",
		Method:      "POST",
		ContentType: "application/json",
	}

	message := interviewApi.GetRequest{
		Ids: postIDs,
	}

	resp := interviewApi.GetResponse{}
	err = endpoint.Send(message, &resp)
	if err != nil {
		return
	}
	interviews = FromApiInterviews(resp.Interviews)
	return
}

func (s *ServiceStorage) SelectInterviewsResults(postIDs []int, userID int) (interviews []interviewModels.InterviewResult, err error) {
	endpoint := service.ServiceEndpoint{
		Service:     s,
		Endpoint:    "/api/interview/list",
		Method:      "POST",
		ContentType: "application/json",
	}

	message := interviewApi.GetUniversalRequest{
		PostIDs:         postIDs,
		UserID:          userID,
		NotPassedResult: true,
	}

	resp := interviewApi.GetUniversalResponse{}
	err = endpoint.Send(message, &resp)
	if err != nil {
		return
	}
	interviews = resp.Interviews
	//interviews = FromApiInterviews(resp.Interviews)
	return
}

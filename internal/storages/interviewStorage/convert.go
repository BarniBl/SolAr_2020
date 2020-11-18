package interviewStorage

import(
	interviewModels "github.com/Solar-2020/Interview-Backend/pkg/models"
	"github.com/Solar-2020/SolAr_Backend_2020/internal/models"
)

func ToApiAnswer(answer models.Answer) interviewModels.Answer {
	return interviewModels.Answer{
		ID:          interviewModels.AnswerID(answer.ID),
		Text:        answer.Text,
		InterviewID: interviewModels.InterviewID(answer.InterviewID),
	}
}

func FromApiAnswer(answer interviewModels.Answer) models.Answer {
	return models.Answer{
		ID:          int(answer.ID),
		Text:        answer.Text,
		InterviewID: int(answer.InterviewID),
	}
}

func ToApiInterviews(src []models.Interview) []interviewModels.Interview {
	return func() []interviewModels.Interview {
		res := make([]interviewModels.Interview, len(src))
		for j, item := range src {
			res[j] = interviewModels.Interview{
				InterviewFrame: interviewModels.InterviewFrame{
					ID:     interviewModels.InterviewID(item.ID),
					Text:   item.Text,
					Type:   interviewModels.InterviewType(item.Type),
					PostID: item.PostID,
					Status: 0,	// TODO: status
				},
				Answers:        func() []interviewModels.Answer{
					r := make([]interviewModels.Answer, len(item.Answers))
					for i, item2 := range item.Answers {
						r[i] = ToApiAnswer(item2)
					}
					return r
				}(),
			}
		}
		return res
	}()
}

func FromApiInterviews(src []interviewModels.Interview) []models.Interview {
	return func() []models.Interview {
		res := make([]models.Interview, len(src))
		for j, item := range src {
			res[j] = models.Interview{
				ID:      int(item.ID),
				Text:    item.Text,
				Type:    int(item.Type),
				PostID:  item.PostID,
				Answers: func() []models.Answer{
					r := make([]models.Answer, len(item.Answers))
					for i, item2 := range item.Answers {
						r[i] = FromApiAnswer(item2)
					}
					return r
				}(),
			}
		}
		return res
	}()
}

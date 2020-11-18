package interviewStorage

import (
	"database/sql"
	interviewModels "github.com/Solar-2020/Interview-Backend/pkg/models"
	"github.com/Solar-2020/SolAr_Backend_2020/internal/models"
	"strconv"
	"strings"
)

type Storage interface {
	InsertInterviews(interviews []models.Interview, postID int) (err error)
	SelectInterviews(postIDs []int) (interviews []models.Interview, err error)
	SelectInterviewsResults(postIDs []int, userID int) (interviews []interviewModels.InterviewResult, err error)
}

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return &storage{
		db: db,
	}
}

func (s *storage) SelectInterviewsResults(postIDs []int, userID int) (interviews []interviewModels.InterviewResult, err error) {
	panic("implement me")
}

func (s *storage) InsertInterviews(interviews []models.Interview, postID int) (err error) {
	if len(interviews) == 0 {
		return
	}

	const sqlQuery = `
	INSERT INTO interviews(text, type, post_id)
	VALUES ($1, $2, $3)
	RETURNING id;`

	tx, err := s.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	for i, _ := range interviews {
		var currentInterviewID int
		err = s.db.QueryRow(sqlQuery, interviews[i].Text, interviews[i].Type, postID).Scan(&currentInterviewID)
		if err != nil {
			return
		}

		err = s.insertAnswers(tx, interviews[i].Answers, currentInterviewID)
		if err != nil {
			return
		}
	}

	err = tx.Commit()

	return
}

func (s *storage) insertAnswers(tx *sql.Tx, answers []models.Answer, interviewID int) (err error) {
	if len(answers) == 0 {
		return
	}

	sqlQueryTemplate := `
	INSERT INTO answers(interview_id, text)
	VALUES `

	if len(answers) == 0 {
		return
	}

	var params []interface{}

	sqlQuery := sqlQueryTemplate + s.createInsertQuery(len(answers), 2)

	for i, _ := range answers {
		params = append(params, interviewID, answers[i].Text)
	}

	for i := 1; i <= len(answers)*2; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	_, err = tx.Exec(sqlQuery, params...)
	return
}

func (s *storage) SelectInterviews(postIDs []int) (interviews []models.Interview, err error) {
	interviews = make([]models.Interview, 0)
	if len(postIDs) == 0 {
		return
	}
	const sqlQueryTemplate = `
	SELECT i.id, i.text, i.type, i.post_id
	FROM interviews AS i
	WHERE i.post_id IN `

	sqlQuery := sqlQueryTemplate + createIN(len(postIDs))

	var params []interface{}

	for i, _ := range postIDs {
		params = append(params, postIDs[i])
	}

	for i := 1; i <= len(postIDs)*1; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	rows, err := s.db.Query(sqlQuery, params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tempInterview models.Interview
		err = rows.Scan(&tempInterview.ID, &tempInterview.Text, &tempInterview.Type, &tempInterview.PostID)
		if err != nil {
			return
		}
		tempInterview.Answers = make([]models.Answer, 0)
		interviews = append(interviews, tempInterview)
	}

	interviewIDs := make([]int, 0)
	for i, _ := range interviews {
		interviewIDs = append(interviewIDs, interviews[i].ID)
	}

	answers, err := s.selectAnswers(interviewIDs)
	if err != nil {
		return
	}

	for _, answer := range answers {
		for i, _ := range interviews {
			if answer.InterviewID == interviews[i].ID {
				interviews[i].Answers = append(interviews[i].Answers, answer)
			}
		}
	}
	return
}

func (s *storage) selectAnswers(interviewIDs []int) (answers []models.Answer, err error) {
	const sqlQueryTemplate = `
	SELECT a.id, a.text, a.interview_id
	FROM answers AS a
	WHERE a.interview_id IN `

	sqlQuery := sqlQueryTemplate + createIN(len(interviewIDs))

	var params []interface{}

	for i, _ := range interviewIDs {
		params = append(params, interviewIDs[i])
	}

	for i := 1; i <= len(interviewIDs)*1; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	rows, err := s.db.Query(sqlQuery, params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tempAnswer models.Answer
		err = rows.Scan(&tempAnswer.ID, &tempAnswer.Text, &tempAnswer.InterviewID)
		if err != nil {
			return
		}
		answers = append(answers, tempAnswer)
	}

	return
}

func createIN(count int) (queryIN string) {
	queryIN = "("
	for i := 0; i < count; i++ {
		queryIN += "?, "
	}
	queryINRune := []rune(queryIN)
	queryIN = string(queryINRune[:len(queryINRune)-2])
	queryIN += ")"
	return
}

func (s *storage) createInsertQuery(sliceLen int, structLen int) (query string) {
	query = ""
	for i := 0; i < sliceLen; i++ {
		query += "("
		for j := 0; j < structLen; j++ {
			query += "?,"
		}
		// delete last comma
		query = strings.TrimRight(query, ",")
		query += "),"
	}
	// delete last comma
	query = strings.TrimRight(query, ",")

	return
}

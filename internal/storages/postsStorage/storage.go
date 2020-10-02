package postsStorage

import (
	"database/sql"
	"github.com/BarniBl/SolAr_2020/internal/models"
	"strconv"
	"strings"
)

const (
	queryReturningID = "RETURNING id;"
)

type Storage interface {
	InsertPost(inputPost models.InputPost) (postID int, err error)
}

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return &storage{
		db: db,
	}
}

func (s *storage) InsertPost(inputPost models.InputPost) (postID int, err error) {
	tx, err := s.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	postID, err = s.insertPost(tx, inputPost)
	if err != nil {
		return
	}

	err = s.insertInterviews(tx, inputPost.Interviews, postID)
	if err != nil {
		return
	}

	err = s.insertPayments(tx, inputPost.Payments, postID)
	if err != nil {
		return
	}

	err = s.insertPhotos(tx, inputPost.Photos, postID)
	if err != nil {
		return
	}

	err = s.insertFiles(tx, inputPost.Files, postID)
	if err != nil {
		return
	}

	err = tx.Commit()

	return
}

func (s *storage) insertPost(tx *sql.Tx, post models.InputPost) (postID int, err error) {
	const sqlQuery = `
	INSERT INTO posts(create_by, publish_date, group_id, text)
	VALUES ($1, $2, $3, $4)
	RETURNING id;`

	err = tx.QueryRow(sqlQuery, post.CreateBy, post.PublishDate, post.GroupID, post.Text).Scan(&postID)

	return
}

func (s *storage) insertInterviews(tx *sql.Tx, interviews []models.Interview, postID int) (err error) {
	const sqlQuery = `
	INSERT INTO interviews(text, type, post_id)
	VALUES ($1, $2, $3)
	RETURNING id;`

	for i, _ := range interviews {
		var currentInterviewID int
		err = tx.QueryRow(sqlQuery, interviews[i].Text, interviews[i].Type, postID).Scan(&currentInterviewID)
		if err != nil {
			return
		}

		err = s.insertAnswers(tx, interviews[i].Answers)
		if err != nil {
			return
		}
	}

	return
}

func (s *storage) insertAnswers(tx *sql.Tx, answers []models.Answer) (err error) {
	sqlQueryTemplate := `
	INSERT INTO interviews(text)
	VALUES `

	var params []interface{}

	sqlQuery := sqlQueryTemplate + s.createInsertQuery(len(answers), 1) + queryReturningID

	for i, _ := range answers {
		params = append(params, answers[i].Text)
	}

	for i := 1; i <= len(answers)*1; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	_, err = tx.Exec(sqlQuery, params)
	return
}

func (s *storage) createInsertQuery(sliceLen int, structLen int) (query string) {
	query = "("
	for i := 0; i < sliceLen; i++ {
		for j := 0; j < structLen; j++ {
			query += "?,"
		}
		query += "),"
	}
	// delete last comma
	query = string([]rune(query)[:len([]rune(query))-1])
	return
}

func (s *storage) insertPayments(tx *sql.Tx, payments []models.Payment, postID int) (err error) {
	sqlQueryTemplate := `
	INSERT INTO payments(post_id, cost, currency_id)
	VALUES `

	var params []interface{}

	sqlQuery := sqlQueryTemplate + s.createInsertQuery(len(payments), 3) + queryReturningID

	for i, _ := range payments {
		params = append(params, postID, payments[i].Cost, payments[i].Currency)
	}

	for i := 1; i <= len(payments)*3; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	_, err = tx.Exec(sqlQuery, params)

	return
}

func (s *storage) insertPhotos(tx *sql.Tx, photos []int, postID int) (err error) {
	sqlQueryTemplate := `
	INSERT INTO photos(post_id, photo_id)
	VALUES `

	var params []interface{}

	sqlQuery := sqlQueryTemplate + s.createInsertQuery(len(photos), 2) + queryReturningID

	for i, _ := range photos {
		params = append(params, postID, photos[i])
	}

	for i := 1; i <= len(photos)*2; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	_, err = tx.Exec(sqlQuery, params)

	return
}

func (s *storage) insertFiles(tx *sql.Tx, files []int, postID int) (err error) {
	sqlQueryTemplate := `
	INSERT INTO files(post_id, file_id)
	VALUES `

	var params []interface{}

	sqlQuery := sqlQueryTemplate + s.createInsertQuery(len(files), 2) + queryReturningID

	for i, _ := range files {
		params = append(params, postID, files[i])
	}

	for i := 1; i <= len(files)*2; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	_, err = tx.Exec(sqlQuery, params)

	return
}

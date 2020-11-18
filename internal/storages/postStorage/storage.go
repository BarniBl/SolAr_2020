package postStorage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Solar-2020/SolAr_Backend_2020/internal/models"
	"strconv"
	"strings"
)

const (
	queryReturningID = "RETURNING id;"
)

type Storage interface {
	InsertPost(inputPost models.InputPost) (postID int, err error)

	UpdatePostStatus(postID int, groupID int, status int) (err error)

	SelectFileIDs(postIDs []int) (matches []models.PostFileMatch, err error)
	SelectPhotoIDs(postIDs []int) (matches []models.PostPhotoMatch, err error)

	SelectPosts(request models.GetPostListRequest) (posts []models.InputPost, err error)
	SelectInterviews(postIDs []int) (interviews []models.Interview, err error)

	SetMark(postID int, mark bool, group int) (err error)
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

func (s *storage) insertPhotos(tx *sql.Tx, photos []int, postID int) (err error) {
	if len(photos) == 0 {
		return
	}

	const sqlQueryTemplate = `
	INSERT INTO photos(post_id, photo_id)
	VALUES `

	var params []interface{}

	sqlQuery := sqlQueryTemplate + s.createInsertQuery(len(photos), 2)

	for i, _ := range photos {
		params = append(params, postID, photos[i])
	}

	for i := 1; i <= len(photos)*2; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	_, err = tx.Exec(sqlQuery, params...)

	return
}

func (s *storage) insertFiles(tx *sql.Tx, files []int, postID int) (err error) {
	if len(files) == 0 {
		return
	}

	const sqlQueryTemplate = `
	INSERT INTO files(post_id, file_id)
	VALUES `

	var params []interface{}

	sqlQuery := sqlQueryTemplate + s.createInsertQuery(len(files), 2)

	for i, _ := range files {
		params = append(params, postID, files[i])
	}

	for i := 1; i <= len(files)*2; i++ {
		sqlQuery = strings.Replace(sqlQuery, "?", "$"+strconv.Itoa(i), 1)
	}

	_, err = tx.Exec(sqlQuery, params...)

	return
}

func (s *storage) UpdatePostStatus(postID int, groupID int, status int) (err error) {
	const sqlQuery = `
	UPDATE posts
	SET status_id = $3
	WHERE id = $1 AND group_id = $2;`

	res, err := s.db.Exec(sqlQuery, postID, groupID, status)
	if err != nil {
		return err
	}

	if changed, _ := res.RowsAffected(); changed != 1 {
		return errors.New("not exist")
	}

	return
}

func (s *storage) SelectPosts(request models.GetPostListRequest) (posts []models.InputPost, err error) {
	posts = make([]models.InputPost, 0)
	sqlQuery := `
	SELECT p.id, p.text, p.group_id, p.publish_date, p.create_by, p.marked
	FROM posts.posts AS p
	WHERE p.group_id = $1
	  AND p.status_id = 2
	  AND p.publish_date <= $2
	  %s
	ORDER BY p.publish_date DESC
	LIMIT $3`

	markCondition := ""
	params := []interface{}{
		request.GroupID, request.StartFrom, request.Limit,
	}
	if request.Mark.Defined {
		markCondition = " AND p.marked = $4"
		params = append(params, request.Mark.Value)
	}

	rows, err := s.db.Query(
		fmt.Sprintf(sqlQuery, markCondition), params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tempPost models.InputPost
		err = rows.Scan(&tempPost.ID, &tempPost.Text, &tempPost.GroupID, &tempPost.PublishDate, &tempPost.CreateBy, &tempPost.Marked)
		if err != nil {
			return
		}
		posts = append(posts, tempPost)
	}

	return
}

func (s *storage) SelectInterviews(postIDs []int) (interviews []models.Interview, err error) {
	const sqlQueryTemplate = `
	SELECT i.id, i.text, i.type, i.post_id
	FROM interview AS i
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
		interviews = append(interviews, tempInterview)
	}

	// TODO SELECT ANSWERS

	return
}

func (s *storage) SelectFileIDs(postIDs []int) (matches []models.PostFileMatch, err error) {
	matches = make([]models.PostFileMatch, 0)
	if len(postIDs) == 0 {
		return
	}
	const sqlQueryTemplate = `
	SELECT f.post_id, f.file_id
	FROM files AS f
	WHERE f.post_id IN `

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
		var tempMatch models.PostFileMatch
		err = rows.Scan(&tempMatch.PostID, &tempMatch.FileID)
		if err != nil {
			return
		}
		matches = append(matches, tempMatch)
	}

	return
}

func (s *storage) SelectPhotoIDs(postIDs []int) (matches []models.PostPhotoMatch, err error) {
	matches = make([]models.PostPhotoMatch, 0)
	if len(postIDs) == 0 {
		return
	}
	const sqlQueryTemplate = `
	SELECT p.post_id, p.photo_id
	FROM photos AS p
	WHERE p.post_id IN `

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
		var tempMatch models.PostPhotoMatch
		err = rows.Scan(&tempMatch.PostID, &tempMatch.PhotoID)
		if err != nil {
			return
		}
		matches = append(matches, tempMatch)
	}

	return
}

func (s *storage) SetMark(postID int, mark bool, group int) (err error) {
	const sqlQueryTemplate = `
	UPDATE posts SET marked=$1 WHERE id=$2 and group_id=$3`
	res, err := s.db.Exec(sqlQueryTemplate, mark, postID, group)
	if err != nil {
		err = errors.New("bad values")
		return
	}

	if affected, err := res.RowsAffected(); err != nil || affected != 1 {
		return errors.New("not changed")
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

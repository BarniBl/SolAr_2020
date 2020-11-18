package models

import (
	account "github.com/Solar-2020/Account-Backend/pkg/models"
	"github.com/Solar-2020/Interview-Backend/pkg/models"
	"github.com/shopspring/decimal"
	"time"
)

type OptBool struct {
	Value   bool
	Defined bool
}

type GetPostListRequest struct {
	UserID    int
	GroupID   int
	Limit     int
	StartFrom time.Time
	Mark      OptBool
}

//type MainPost struct {
//	ID          int         `json:"id"`
//	CreateBy    int         `json:"-"`
//	CreatAt     time.Time   `json:"-"`
//	PublishDate time.Time   `json:"publishDate"`
//	GroupID     int         `json:"groupID"`
//	Text        string      `json:"text"`
//	Status      string      `json:"Status"`
//	Interviews  []Interview `json:"interviews"`
//	Payments    []Payment   `json:"payments"`
//}

type InputPost struct {
	ID          int         `json:"id"`
	CreateBy    int         `json:"createBy"`
	CreatAt     time.Time   `json:"-"`
	PublishDate time.Time   `json:"publishDate"`
	GroupID     int         `json:"groupID"`
	Text        string      `json:"text"`
	Status      string      `json:"Status"`
	Photos      []int       `json:"photos"`
	Files       []int       `json:"files"`
	Interviews  []Interview `json:"interviews"`
	Payments    []Payment   `json:"payments"`
	Marked      bool        `json:"marked"`
}

type Post struct {
	ID          int         `json:"id"`
	CreateBy    int         `json:"-"`
	CreatAt     time.Time   `json:"-"`
	PublishDate time.Time   `json:"publishDate"`
	GroupID     int         `json:"groupID"`
	Text        string      `json:"text"`
	Status      string      `json:"Status"`
	Photos      []Photo     `json:"photos"`
	Files       []File      `json:"files"`
	Interviews  []Interview `json:"interviews"`
	Payments    []Payment   `json:"payments"`
	Order       int         `json:"-"`
	Marked      bool        `json:"marked"`
}

type CreateRequest struct {
	CreateBy int       `json:"createBy"`
	GroupID  int       `json:"groupID"`
	PostID   int       `json:"postID"`
	Payments []Payment `json:"payments"`
}

type Payment struct {
	ID             int             `json:"id"`
	GroupID        int             `json:"groupID"`
	PostID         int             `json:"postID"`
	CreateBy       int             `json:"createBy"`
	TotalCost      decimal.Decimal `json:"totalCost"`
	PaymentAccount string          `json:"paymentAccount"`
}

type PostResult struct {
	ID          int                      `json:"id"`
	Author      account.User             `json:"author"`
	CreateBy    int                      `json:"-"`
	CreatAt     time.Time                `json:"-"`
	PublishDate time.Time                `json:"publishDate"`
	GroupID     int                      `json:"groupID"`
	Text        string                   `json:"text"`
	Status      string                   `json:"Status"`
	Photos      []Photo                  `json:"photos"`
	Files       []File                   `json:"files"`
	Interviews  []models.InterviewResult `json:"interviews"`
	Payments    []Payment                `json:"payments"`
	Order       int                      `json:"-"`
	Marked      bool                     `json:"marked"`
}

type Posts struct {
	Posts []PostResult
}

func (p *Posts) Len() int {
	return len(p.Posts)
}

func (p *Posts) Swap(i, j int) {
	p.Posts[i], p.Posts[j] = p.Posts[j], p.Posts[i]
}

func (p *Posts) Less(i, j int) bool {
	return p.Posts[i].Order < p.Posts[j].Order
}

type UserRequest struct {
	UserID int
}

type MarkPost struct {
	UserRequest
	PostID  int
	GroupID int
	Mark    bool
}

type DeletePostRequest struct {
	UserRequest
	PostID  int `json:"postId"`
	GroupID int `json:"groupId"`
}

type Interview struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Type    int      `json:"type"`
	PostID  int      `json:"postID"`
	Answers []Answer `json:"answers"`
}

type Answer struct {
	ID          int    `json:"id"`
	Text        string `json:"text"`
	InterviewID int    `json:"interviewID"`
}

type AclAction int

const (
	ActionGetList AclAction = iota
	ActionCreatePost
)

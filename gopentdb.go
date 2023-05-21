package gopentdb

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Categories struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type categoryQuestionCountRes struct {
	TotalQuestionCount       int64 `json:"total_question_count"`
	TotalEasyQuestionCount   int64 `json:"total_easy_question_count"`
	TotalMediumQuestionCount int64 `json:"total_medium_question_count"`
	TotalHardQuestionCount   int64 `json:"total_hard_question_count"`
}

type categoryCountRes struct {
	CategoryId            int64                    `json:"category_id"`
	CategoryQuestionCount categoryQuestionCountRes `json:"category_question_count"`
}

type ResponseCode int

const (
	Success          ResponseCode = 0
	NoResults        ResponseCode = 1
	InvalidParameter ResponseCode = 2
	TokenNotFound    ResponseCode = 3
	TokenEmpty       ResponseCode = 4
)

type categoriesRes struct {
	TriviaCategories []Categories `json:"trivia_categories"`
}

type sessionTokenRes struct {
	ResponseCode    ResponseCode `json:"response_code"`
	ResponseMessage string       `json:"response_message"`
	Token           string       `json:"token"`
}

type Config struct {
	BaseUrl string
}

type OpenTDB struct {
	client  *resty.Client
	baseUrl string
}

func New(c Config) OpenTDB {
	o := OpenTDB{}
	o.init(c.BaseUrl)
	return o
}

func (o *OpenTDB) init(baseUrl string) {
	o.baseUrl = baseUrl
	o.client = resty.New()
	o.client.SetBaseURL(o.baseUrl)
}

// Check if OpenTDB is Up and running.
func (o *OpenTDB) Ping() (bool, error) {
	resp, err := o.client.R().Get("/")
	if err != nil {
		return false, err
	}
	if resp.StatusCode() == 200 {
		return true, nil
	} else {
		return false, nil
	}
}

// Get List of Categories from OpenTDB
func (o *OpenTDB) GetCategories() ([]Categories, error) {
	resp, err := o.client.R().
		SetHeader("Accept", "application/json").
		SetResult(&categoriesRes{}).
		Get("/api_category.php")
	if err != nil {
		return make([]Categories, 0), err
	}
	categoriesList := resp.Result().(*categoriesRes).TriviaCategories
	return categoriesList, nil
}

// Retrieve a Session Token
func (o *OpenTDB) GetSessionToken() (string, error) {
	resp, err := o.client.R().
		SetHeader("Accept", "application/json").
		SetResult(&sessionTokenRes{}).
		Get("/api_token.php?command=request")
	if err != nil {
		return "", err
	}
	token := resp.Result().(*sessionTokenRes).Token
	return token, nil
}

// Gets Question Count for a given category Id
func (o *OpenTDB) GetCategoryCount(category int64) (int64, error) {
	resp, err := o.client.R().
		SetHeader("Accept", "application/json").
		SetResult(&categoryCountRes{}).
		SetQueryParam("category", fmt.Sprint(category)).
		Get("/api_count.php")
	if err != nil {
		return -1, err
	}
	count := resp.Result().(*categoryCountRes).CategoryQuestionCount.TotalQuestionCount
	return count, nil
}

type QuestionDifficulty string

const (
	Easy   QuestionDifficulty = "easy"
	Medium QuestionDifficulty = "medium"
	Hard   QuestionDifficulty = "hard"
)

type QuestionType string

const (
	Multiple QuestionType = "multiple"
	Boolean  QuestionType = "boolean"
)

type Encoding string

const (
	Default     Encoding = ""
	LegacyUrl   Encoding = "urlLegacy"
	UrlEncoding Encoding = "url3986"
	Base64      Encoding = "base64"
)

type Question struct {
	CategoryName       string             `json:"category"`
	QuestionType       QuestionType       `json:"type"`
	QuestionDifficulty QuestionDifficulty `json:"difficulty"`
	Question           string             `json:"question"`
	CorrectAnswer      string             `json:"correct_answer"`
	IncorrectAnswers   []string           `json:"incorrect_answers"`
}

type questionRes struct {
	ResponseCode ResponseCode `json:"response_code"`
	Results      []Question   `json:"results"`
}

type QuestionParams struct {
	Amount             int64              `default:"10"`
	Category           int64              `default:"-1"`
	QuestionDifficulty QuestionDifficulty `default:""`
	QuestionType       QuestionType       `default:""`
	Encoding           `default:""`
	Token              string `default:""`
}

// Get Questions based on Amount, Category, Difficulty and Type
func (o *OpenTDB) GetQuestions(params QuestionParams) ([]Question, error) {
	paramsMap := make(map[string]string)
	if params.Amount > 50 {
		paramsMap["amount"] = "50"
	} else if params.Amount > 0 {
		paramsMap["amount"] = fmt.Sprint(params.Amount)
	} else {
		paramsMap["amount"] = "10"
	}
	if params.Category > 0 {
		paramsMap["category"] = fmt.Sprint(params.Category)
	}
	if len(params.QuestionDifficulty) > 0 {
		paramsMap["difficulty"] = fmt.Sprint(params.QuestionDifficulty)
	}
	if len(params.QuestionType) > 0 {
		paramsMap["type"] = fmt.Sprint(params.QuestionType)
	}
	if len(params.Encoding) > 0 {
		paramsMap["encoding"] = fmt.Sprint(params.Encoding)
	}
	if len(params.Token) > 0 {
		paramsMap["token"] = fmt.Sprint(params.Token)
	}
	resp, err := o.client.R().
		SetHeader("Accept", "application/json").
		SetResult(&questionRes{}).
		SetQueryParams(paramsMap).
		Get("/api.php")
	if err != nil {
		return make([]Question, 0), err
	}
	questions := resp.Result().(*questionRes).Results
	return questions, nil
}

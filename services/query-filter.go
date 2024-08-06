package services

import (
	"regexp"
	"strconv"
	"strings"
)

type QueryFilter struct {
	Page     int            `json:"page"`
	Limit    int            `json:"limit"`
	Skip     int            `json:"skip"`
	Select   []string       `json:"select"`
	Search   []string       `json:"search"`
	Where    map[string]any `json:"where"`
	Sort     string         `json:"sort"`
	Populate []string       `json:"populate"`
}

func NewQueryFilter(query map[string][]string) (QueryFilter, error) {
	queryFilter := QueryFilter{
		Page:     1,
		Limit:    0,
		Skip:     0,
		Select:   []string{},
		Search:   []string{},
		Where:    make(map[string]any),
		Sort:     "id asc",
		Populate: []string{},
	}

	if err := queryFilter.initQuery(query); err != nil {
		return QueryFilter{}, err
	}

	return queryFilter, nil
}

func (q *QueryFilter) initQuery(query map[string][]string) error {
	var err error

	matchWhere, _ := regexp.Compile(`where\[(.*)\]`)
	matchSearch, _ := regexp.Compile(`search\[(.*)\]`)

	for key, value := range query {
		if key == "page" {
			q.Page, err = strconv.Atoi(value[0])
		} else if key == "limit" {
			q.Limit, err = strconv.Atoi(value[0])
		} else if key == "sort" {
			q.Sort = value[0]
		} else if key == "select" {
			q.Select = strings.Split(value[0], ",")
		} else if key == "populate" {
			q.Populate = strings.Split(value[0], ",")
		} else if matchWhere.MatchString(key) {
			matchKey := matchWhere.ReplaceAllString(key, "$1")
			q.Where[matchKey] = strings.Split(value[0], ",")
		} else if matchSearch.MatchString(key) {
			matchKey := matchSearch.ReplaceAllString(key, "$1")
			values := strings.Split(value[0], ",")
			for _, value := range values {
				expre := matchKey + " LIKE '%" + value + "%'"
				q.Search = append(q.Search, expre)
			}
		}

		if err != nil {
			return err
		}
	}

	q.Skip = (q.Page - 1) * q.Limit

	return nil
}

func (q *QueryFilter) GetPage() int {
	return q.Page
}

func (q *QueryFilter) GetLimit() int {
	return q.Limit
}

func (q *QueryFilter) GetWhere() map[string]any {
	return q.Where
}

func (q *QueryFilter) GetSkip() int {
	return q.Skip
}

func (q *QueryFilter) GetSort() string {
	return q.Sort
}

func (q *QueryFilter) GetSelect() []string {
	return q.Select
}

func (q *QueryFilter) GetPopulate() []string {
	return q.Populate
}

func (q *QueryFilter) GetSearch() string {
	return strings.Join(q.Search, " AND ")
}

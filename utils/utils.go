package utils

import (
	"errors"
	"strconv"
	"strings"
)

type Pagination struct {
	Start int
	End   int
}

var InvalidPaginationParamsError = errors.New("Invalid pagination params")

func FormatedPagination(pagination string) (*Pagination, error) {
	str := strings.Split(pagination[1:len(pagination)-1], ",")
	res := make([]int, 2)
	for i, s := range str {
		num, err := strconv.Atoi(s)
		if err != nil {
			return &Pagination{}, err
		}
		res[i] = num
	}

	if res[0] < 0 || res[1] < 0 {
		return &Pagination{}, InvalidPaginationParamsError
	}
	return &Pagination{Start: res[0], End: res[1]}, nil
}

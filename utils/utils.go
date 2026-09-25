package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

type Pagination struct {
	Start int
	End   int
}

var UniqueViolationErrorCode = "23505"
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

func IsUniqueViolationError(err error) bool {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == UniqueViolationErrorCode {
		return true
	}
	return false
}

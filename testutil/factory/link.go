package factory

import (
	"context"
	"github.com/seeu359/go-project-278/links"
	"testing"
)

func CreateLinks(t *testing.T, q *links.Queries) error {
	t.Helper()
	ctx := context.Background()
	links := []links.CreateLinkParams{
		{Url: "https://ya.ru", ShortName: "short"},
		{Url: "https://yandex.ru", ShortName: "test"},
	}

	for _, l := range links {
		_, err := q.CreateLink(ctx, l)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateLink(t *testing.T, q *links.Queries) (int64, error) {
	t.Helper()
	ctx := context.Background()
	link := links.CreateLinkParams{Url: "https://ya.ru", ShortName: "short"}
	id, err := q.CreateLink(ctx, links.CreateLinkParams{Url: link.Url, ShortName: link.ShortName})
	if err != nil {
		return 0, err
	}
	return id, nil
}

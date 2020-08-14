package trello

import (
	"fmt"
)

var (
	BoardsEndpoint     = func() string { return buildURL("members/me/boards") }
	BoardEndpoint      = func(boardID string) string { return buildURL(fmt.Sprintf("boards/%s", boardID)) }
	BoardListsEndpoint = func(boardID string) string { return buildURL(fmt.Sprintf("boards/%s/lists", boardID)) }
	CreateListEndpoint = func(boardID string) string { return buildURL(fmt.Sprintf("boards/%s/lists", boardID)) }
)

const baseURL = "https://api.trello.com/1"

func buildURL(path string) string {
	return fmt.Sprintf("%s/%s", baseURL, path)
}

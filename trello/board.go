package trello

import "fmt"

type Board struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Desc            string        `json:"desc"`
	DescData        string        `json:"descData"`
	Closed          bool          `json:"closed"`
	IDOrganization  string        `json:"idOrganization"`
	IDEnterprise    string        `json:"idEnterprise"`
	IDBoardSource   string        `json:"idBoardSource"`
	PremiumFeatures []interface{} `json:"premiumFeatures"`
	Pinned          bool          `json:"pinned"`
	URL             string        `json:"url"`
	ShortURL        string        `json:"shortUrl"`
	Prefs           struct {
		PermissionLevel       string `json:"permissionLevel"`
		HideVotes             bool   `json:"hideVotes"`
		Voting                string `json:"voting"`
		Comments              string `json:"comments"`
		SelfJoin              bool   `json:"selfJoin"`
		CardCovers            bool   `json:"cardCovers"`
		IsTemplate            bool   `json:"isTemplate"`
		CardAging             string `json:"cardAging"`
		CalendarFeedEnabled   bool   `json:"calendarFeedEnabled"`
		Background            string `json:"background"`
		BackgroundImage       string `json:"backgroundImage"`
		BackgroundImageScaled []struct {
			Width  int    `json:"width"`
			Height int    `json:"height"`
			URL    string `json:"url"`
		} `json:"backgroundImageScaled"`
	} `json:"prefs"`
	LabelNames struct {
		Green  string `json:"green"`
		Yellow string `json:"yellow"`
		Orange string `json:"orange"`
		Red    string `json:"red"`
		Purple string `json:"purple"`
		Blue   string `json:"blue"`
		Sky    string `json:"sky"`
		Lime   string `json:"lime"`
		Pink   string `json:"pink"`
		Black  string `json:"black"`
	} `json:"labelNames"`
	Limits struct {
		Attachments struct {
			PerBoard struct {
				Status    string `json:"status"`
				DisableAt int    `json:"disableAt"`
				WarnAt    int    `json:"warnAt"`
			} `json:"perBoard"`
		} `json:"attachments"`
	} `json:"limits"`
	Starred     bool `json:"starred"`
	Memberships []struct {
		ID          string `json:"id"`
		IDMember    string `json:"idMember"`
		MemberType  string `json:"memberType"`
		Unconfirmed bool   `json:"unconfirmed"`
		Deactivated bool   `json:"deactivated"`
	} `json:"memberships"`
	ShortLink         string        `json:"shortLink"`
	Subscribed        bool          `json:"subscribed"`
	PowerUps          []interface{} `json:"powerUps"`
	DateLastActivity  string        `json:"dateLastActivity"`
	DateLastView      string        `json:"dateLastView"`
	IDTags            []interface{} `json:"idTags"`
	DatePluginDisable string        `json:"datePluginDisable"`
	CreationMethod    string        `json:"creationMethod"`
	IxUpdate          int           `json:"ixUpdate"`
	TemplateGallery   string        `json:"templateGallery"`
	EnterpriseOwned   bool          `json:"enterpriseOwned"`
}

type List struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	IDBoard    string      `json:"idBoard"`
	Closed     bool        `json:"closed"`
	Subscribed bool        `json:"subscribed"`
	Pos        int         `json:"pos"`
	SoftLimit  interface{} `json:"softLimit"`
}

func (t *Trello) Boards() ([]*Board, error) {
	url := BoardsEndpoint()
	var boards []*Board
	if err := t.handler.getJSON(url, &boards); err != nil {
		return nil, err
	}

	return boards, nil
}

func (t *Trello) SearchBoardByName(boardName string) (*Board, error) {
	boards, err := t.Boards()
	if err != nil {
		return nil, err
	}
	for _, b := range boards {
		if b.Name == boardName {
			return b, nil
		}
	}
	return nil, fmt.Errorf("Can't find board %s", boardName)
}

func (t *Trello) SearchBoardByID(boardID string) (*Board, error) {
	url := BoardEndpoint(boardID)
	var board *Board
	if err := t.handler.getJSON(url, &board); err != nil {
		return nil, err
	}

	return board, nil
}

func (t *Trello) BoardLists(boardID string) ([]*List, error) {
	url := BoardListsEndpoint(boardID)
	var lists []*List
	if err := t.handler.getJSON(url, &lists); err != nil {
		return nil, err
	}

	return lists, nil
}

type ErrListNotFound string

func (e ErrListNotFound) Error() string {
	return fmt.Sprintf("Cannot find list %s", string(e))
}

func (t *Trello) SearchListByName(boardID, listName string) (*List, error) {
	lists, err := t.BoardLists(boardID)
	if err != nil {
		return nil, err
	}
	for _, l := range lists {
		if l.Name == listName {
			return l, nil
		}
	}
	return nil, ErrListNotFound(listName)
}

func (t *Trello) CreateList(boardID, listName string) (*List, error) {
	url := CreateListEndpoint(boardID)

	var list *List
	args := map[string]string{
		"name": listName,
	}

	if err := t.handler.postJSON(url, args, &list); err != nil {
		return nil, err
	}

	return list, nil
}

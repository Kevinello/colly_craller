package storage

type Category struct {
	DistinctID string `json:"distinctId"`
	DataSetID  int    `json:"dataSetId"`
	Title1     string `json:"title1"`
	Title2     string `json:"title2"`
	Title3     string `json:"title3"`
	Action1    string `json:"action1"`
	ID         string `json:"id"`
	Action2    string `json:"action2"`
	Action3    string `json:"action3"`
	Pos        int    `json:"__pos__"`
	Track      string `json:"__track__"`
	ContentID  string `json:"contentId,omitempty"`
}

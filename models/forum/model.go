package forum

type Forum struct {
	Posts   int64 `json:"posts,omitempty"`
	Slug    string
	Threads int32 `json:"threads,omitempty"`
	Title   string
	User    string
}

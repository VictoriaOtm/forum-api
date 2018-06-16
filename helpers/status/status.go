package status

//go:generate easyjson -snake_case

import "github.com/VictoriaOtm/forum-api/database"

//easyjson:json
type Status struct {
	Forum  int32
	Post   int64
	Thread int32
	User   int32
}

const query = `SELECT 
(SELECT count(*) FROM t_forum), 
(SELECT count(*) FROM t_posts), 
(SELECT count(*) FROM t_thread), 
(SELECT count(*) FROM t_user)`

func (s *Status) Get() error {
	err := database.DB.QueryRow(query).
		Scan(&s.Forum, &s.Post, &s.Thread, &s.User)
	return err
}

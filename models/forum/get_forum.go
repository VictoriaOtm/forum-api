package forum

import "github.com/VictoriaOtm/forum-api/database"

func (f *Forum) Get(slug string) error {
	err := database.DBConnPool.QueryRow(database.GetForumBySlug, slug).
		Scan(&f.Slug, &f.Title, &f.User, &f.Threads, &f.Posts)
	return err
}

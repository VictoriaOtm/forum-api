package forum

import "github.com/VictoriaOtm/forum-api/database"

func (f *Forum) Create() error {
	tx, err := database.DBConnPool.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	err = tx.QueryRow(database.GetForumBySlug, f.Slug).
		Scan(&f.Slug, &f.Title, &f.User, &f.Threads, &f.Posts)
	if err == nil {
		return database.ErrorForumConflict
	}

	err = tx.QueryRow(database.GetRealUserNickname, f.User).Scan(&f.User)
	if err != nil {
		return database.ErrorUserNotExists
	}

	_, err = tx.Exec(database.CreateForum, f.Slug, f.Title, f.User)
	if err != nil {
		return err
	}

	return nil
}

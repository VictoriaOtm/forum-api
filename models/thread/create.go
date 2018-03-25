package thread

import (
	"github.com/VictoriaOtm/forum-api/database"
	model "github.com/VictoriaOtm/forum-api/models/user"
)

func (t *Thread) Create() error {
	tx, err := database.DBConnPool.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	if t.Slug != nil {
		err = tx.QueryRow(database.GetThreadBySlug, t.Slug).
			Scan(&t.Id, &t.Slug, &t.Author, &t.Created, &t.Forum, &t.Message, &t.Title, &t.Votes)
		if err == nil {
			return database.ErrorThreadConflict
		}
	}

	user := model.User{}
	err = tx.QueryRow("GetUserByNickname", t.Author).Scan(&user.Nickname, &user.Email, &user.Fullname, &user.About)
	if err != nil {
		return database.ErrorUserNotExists
	}
	t.Author = user.Nickname

	err = tx.QueryRow(database.GetRealForumSlug, t.Forum).Scan(&t.Forum)
	if err != nil {
		return database.ErrorForumNotExists
	}

	err = tx.QueryRow(database.CreateThread, t.Slug, t.Author, t.Created, t.Forum, t.Message, t.Title).
		Scan(&t.Id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("CreateForumUsers", t.Forum, user.Nickname, user.Email, user.Fullname, user.About)
	if err != nil {
		return err
	}
	return nil
}

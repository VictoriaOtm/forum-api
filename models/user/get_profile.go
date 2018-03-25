package user

import (
	"github.com/VictoriaOtm/forum-api/database"
	"bytes"
	"github.com/jackc/pgx"
)

func (usr *User) GetProfile() error {
	return database.DBConnPool.QueryRow(database.GetUserByNickname, usr.Nickname).
		Scan(&usr.Nickname, &usr.Email, &usr.Fullname, &usr.About)
}

func (usrArr *Arr) GetForumUsers(forumSlug string, limit []byte, since []byte, desc []byte) error {
	err := database.DBConnPool.QueryRow(database.GetRealForumSlug, forumSlug).Scan(&forumSlug)
	if err != nil {
		return database.ErrorForumNotExists
	}
	descBool := bytes.Equal([]byte("true"), desc)

	var rows *pgx.Rows
	if since == nil {
		if !descBool {
			rows, err = database.DBConnPool.Query(database.GetForumUsersNsNd, forumSlug, limit)
		} else {
			rows, err = database.DBConnPool.Query(database.GetForumUsersNsYd, forumSlug, limit)
		}
	} else {
		if !descBool {
			rows, err = database.DBConnPool.Query(database.GetForumUsersYsNd, forumSlug, since, limit)
		} else {
			rows, err = database.DBConnPool.Query(database.GetForumUsersYsYd, forumSlug, since, limit)
		}
	}

	if err != nil {
		return err
	}

	l := len(*usrArr)
	for rows.Next() {
		l += 1
		*usrArr = (*usrArr)[:l]

		err = rows.Scan(&(*usrArr)[l-1].Nickname, &(*usrArr)[l-1].Email,
			&(*usrArr)[l-1].Fullname, &(*usrArr)[l-1].About)
		if err != nil {
			return err
		}
	}

	return nil
}

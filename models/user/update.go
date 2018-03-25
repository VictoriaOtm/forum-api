package user

import "github.com/VictoriaOtm/forum-api/database"

func (u *User) Update(upd *Update) error {
	tx, err := database.DBConnPool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if upd.Email != nil {
		err = tx.QueryRow(database.GetUserByEmail, upd.Email).
			Scan(&u.Nickname, &u.Email, &u.About, &u.Fullname)

		if err == nil {
			return database.ErrorUserConflict
		}
	}

	err = tx.QueryRow(database.ReplaceUserInfo, upd.Email, upd.Fullname, upd.About, u.Nickname).
		Scan(&u.Nickname, &u.Email, &u.Fullname, &u.About)
	if err != nil {
		return database.ErrorUserNotExists
	}

	tx.Commit()
	return nil
}

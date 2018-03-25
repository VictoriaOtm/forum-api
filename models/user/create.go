package user

import "github.com/VictoriaOtm/forum-api/database"

func (u *User) Create() (err error, duplicateUsers Arr) {
	tx, err := database.DBConnPool.Begin()
	if err != nil {
		return err, nil
	}
	defer tx.Commit()

	r, err := tx.Exec(database.CreateUser, u.Nickname, u.Email, u.Fullname, u.About)
	if err == nil && r.RowsAffected() == 1 {
		return nil, nil
	}

	rows, err := tx.Query(database.GetUserByNicknameOrEmail, u.Nickname, u.Email)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	for rows.Next() {
		usr := User{}

		err = rows.Scan(&usr.Nickname, &usr.Email, &usr.Fullname, &usr.About)
		if err != nil {
			return err, nil
		}

		duplicateUsers = append(duplicateUsers, usr)
	}

	return database.ErrorUserConflict, duplicateUsers
}

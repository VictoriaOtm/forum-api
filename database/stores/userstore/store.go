package userstore

import (
	"context"

	"strings"

	"github.com/VictoriaOtm/forum-api/database"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/jackc/pgx"
)

func PrepareStatements() {
	database.DB.Prepare("getUserByNickname", getUserByNickname)

	database.DB.Prepare("getForumUsersNsNd", getForumUsersNsNd)
	database.DB.Prepare("getForumUsersNsYd", getForumUsersNsYd)
	database.DB.Prepare("getForumUsersYsNd", getForumUsersYsNd)
	database.DB.Prepare("getForumUsersYsYd", getForumUsersYsYd)

	database.DB.Prepare("createForumUsers", createForumUsers)
}

const getUserByNickname = `
SELECT nickname::TEXT, email::TEXT, fullname, about 
FROM t_user 
WHERE nickname=$1`

func (u *User) Get(username interface{}) error {
	return database.DB.QueryRow("getUserByNickname", username).
		Scan(&u.Nickname, &u.Email, &u.Fullname, &u.About)
}

const getUserByNicknameOrEmail = `
SELECT nickname::TEXT, email::TEXT, fullname, about 
FROM t_user 
WHERE nickname=$1 or email=$2`

func (*User) GetByNicknameOrEmail(nickname, email string) (UserSlice, error) {
	rows, err := database.DB.Query(getUserByNicknameOrEmail, nickname, email)
	if err != nil {
		return nil, err
	}
	us := UserSlice{}
	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.Nickname, &u.Email, &u.Fullname, &u.About)
		if err != nil {
			return nil, err
		}
		us = append(us, u)
	}

	return us, nil
}

func GetByNicknames(nicknames *treemap.Map) error {
	batch := database.DB.BeginBatch()
	defer batch.Close()

	for _, nick := range nicknames.Keys() {
		batch.Queue("getUserByNickname", []interface{}{nick}, nil, nil)
	}

	err := batch.Send(context.Background(), nil)
	if err != nil {
		return err
	}

	for i := 0; i < nicknames.Size(); i++ {
		u := User{}
		err = batch.QueryRowResults().Scan(
			&u.Nickname,
			&u.Email,
			&u.Fullname,
			&u.About,
		)

		if err != nil {
			return err
		}
		nicknames.Put(strings.ToLower(u.Nickname), u)
	}

	return nil
}

const createForumUsers = `
INSERT INTO t_forum_user(slug, nickname, email, fullname, about)
VALUES($1, $2, $3, $4, $5)
ON CONFLICT DO NOTHING`

func StoreInForumUserTable(forumSlug string, users []interface{}) {
	batch := database.DB.BeginBatch()
	defer batch.Close()

	for _, u := range users {
		usr := u.(User)
		batch.Queue("createForumUsers", []interface{}{
			forumSlug,
			usr.Nickname,
			usr.Email,
			usr.Fullname,
			usr.About,
		}, nil, nil)
	}

	batch.Send(context.Background(), nil)
}

const createUser = `
INSERT INTO t_user(nickname, email, fullname, about) 
VALUES ($1, $2, $3, $4)`

func (u *User) Insert() error {
	tx := database.TxMustBegin()
	defer tx.Commit()

	_, err := tx.Exec(createUser, u.Nickname, u.Email, u.Fullname, u.About)
	if err != nil {
		tx.Rollback()
	}
	return err
}

const (
	getForumUsersNsNd = `
SELECT nickname::TEXT, email::TEXT, fullname, about
FROM t_forum_user
WHERE slug=$1
ORDER BY nickname::citext
LIMIT $2::TEXT::INT`

	getForumUsersNsYd = `
SELECT nickname::TEXT, email::TEXT, fullname, about
FROM t_forum_user
WHERE slug=$1
ORDER BY nickname::citext DESC
LIMIT $2::TEXT::INT`

	getForumUsersYsNd = `
SELECT nickname::TEXT, email::TEXT, fullname, about
FROM t_forum_user
WHERE slug=$1 AND nickname>$2::TEXT::citext
ORDER BY nickname::citext
LIMIT $3::TEXT::INT`

	getForumUsersYsYd = `
SELECT nickname::TEXT, email::TEXT, fullname, about
FROM t_forum_user
WHERE slug=$1 AND nickname<$2::TEXT::citext
ORDER BY nickname::citext DESC
LIMIT $3::TEXT::INT`
)

func (us *UserSlice) Get(
	forumSlug interface{},
	limit []byte,
	since []byte,
	desc bool,
) error {
	var rows *pgx.Rows
	var err error

	if since == nil {
		if !desc {
			rows, err = database.DB.Query("getForumUsersNsNd", forumSlug, limit)
		} else {
			rows, err = database.DB.Query("getForumUsersNsYd", forumSlug, limit)
		}
	} else {
		if !desc {
			rows, err = database.DB.Query("getForumUsersYsNd", forumSlug, since, limit)
		} else {
			rows, err = database.DB.Query("getForumUsersYsYd", forumSlug, since, limit)
		}
	}

	if err != nil {
		return err
	}

	*us = (*us)[:100]
	l := 0
	for rows.Next() {
		err = rows.Scan(
			&(*us)[l].Nickname,
			&(*us)[l].Email,
			&(*us)[l].Fullname,
			&(*us)[l].About,
		)

		if err != nil {
			return err
		}

		l += 1
	}

	*us = (*us)[:l]

	return nil
}

const ReplaceUserInfo = `
UPDATE t_user 
SET email=COALESCE($1, email),
	fullname=COALESCE($2, fullname),
	about=COALESCE($3, about)
WHERE nickname=$4
RETURNING nickname::TEXT, email::TEXT, fullname, about`

func (u *User) Update(upd UserUpdate) error {
	tx := database.TxMustBegin()
	defer tx.Commit()

	err := tx.QueryRow(ReplaceUserInfo, upd.Email, upd.Fullname, upd.About, u.Nickname).
		Scan(&u.Nickname, &u.Email, &u.Fullname, &u.About)
	if err != nil {
		tx.Rollback()
	}
	return err
}

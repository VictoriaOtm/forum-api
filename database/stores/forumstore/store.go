package forumstore

import (
	"log"

	"github.com/VictoriaOtm/forum-api/database"
)

func PrepareStatements() {
	database.DB.Prepare("GetForumBySlug", GetForumBySlug)
}

const GetForumBySlug = `
SELECT slug::TEXT, title, f_user, threads, posts
FROM t_forum
WHERE slug=$1`

func (v *Forum) Get(slug interface{}) error {
	return database.DB.QueryRow("GetForumBySlug", slug).
		Scan(&v.Slug, &v.Title, &v.User, &v.Threads, &v.Posts)
}

const createForum = `
INSERT INTO t_forum(slug, title, f_user)
VALUES($1, $2, $3)
`

func (v *Forum) Insert() error {
	tx := database.TxMustBegin()
	defer tx.Commit()

	_, err := tx.Exec(createForum, v.Slug, v.Title, v.User)
	if err != nil {
		err2 := tx.Rollback()
		log.Println(err2)
	}
	return err
}

const updateForumPostCount = `
UPDATE t_forum
SET posts=posts+$2
WHERE slug=$1`

func UpdateForumPosts(slug string, count int64) error {
	tx := database.TxMustBegin()
	defer tx.Commit()

	_, err := tx.Exec(updateForumPostCount, slug, count)
	if err != nil {
		tx.Rollback()
	}
	return nil
}

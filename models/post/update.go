package post

import (
	"github.com/VictoriaOtm/forum-api/database"
	"log"
)

const q = `
UPDATE t_posts
SET is_edited = (SELECT (($1::TEXT IS NOT NULL OR is_edited) AND $1::TEXT != message)),
message = COALESCE($1, message)
WHERE id = $2
RETURNING message, is_edited`

func (p *Post) Update(upd PostUpdate, id string) error {
	tx, _ := database.DBConnPool.Begin()
	defer tx.Commit()

	err := tx.QueryRow(database.GetPost, id).
		Scan(&p.Id, &p.Author, &p.Forum, &p.Thread,
		&p.Created, &p.IsEdited, &p.Parent, &p.Message)
	if err != nil {
		return err
	}

	err = tx.QueryRow(q, upd.Message, p.Id).Scan(&p.Message, &p.IsEdited)
	if err != nil {
		log.Println(err)
		tx.Rollback()
	}

	return err
}

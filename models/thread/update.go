package thread

import "github.com/VictoriaOtm/forum-api/database"

const qupd = `
UPDATE t_thread
SET message = COALESCE($1, message),
  title     = COALESCE($2, title)
WHERE id = $3
RETURNING message, title`

func (t *Thread) Update(upd Update) error {
	tx, err := database.DBConnPool.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	err = tx.QueryRow(qupd, upd.Message, upd.Title, t.Id).Scan(&t.Message, &t.Title)
	if err != nil {
		tx.Rollback()
	}

	return err
}

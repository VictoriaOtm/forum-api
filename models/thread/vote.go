package thread

import "github.com/VictoriaOtm/forum-api/database"

const q = `
WITH find_diff AS (SELECT ($3 - COALESCE((SELECT vote
                                          FROM t_vote
                                          WHERE nickname = $2 and thread_id = $1), 0)) as diff),
    update_vote AS (INSERT INTO t_vote VALUES ($1, $2, $3)
  ON CONFLICT ON CONSTRAINT t_vote_thread_id_nickname_idx
    DO UPDATE SET vote = $3)
UPDATE t_thread
SET votes = votes + (SELECT diff
                     FROM find_diff)
WHERE id = $1
RETURNING votes`

func (t *Thread) PutVote(v Vote) error {
	tx, err := database.DBConnPool.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	err = tx.QueryRow(database.GetRealUserNickname, v.Nickname).Scan(&v.Nickname)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.QueryRow(q, t.Id, v.Nickname, v.Voice).Scan(&t.Votes)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

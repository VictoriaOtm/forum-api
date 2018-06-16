package threadstore

import (
	"github.com/VictoriaOtm/forum-api/database"
	"github.com/VictoriaOtm/forum-api/helpers"
	"github.com/jackc/pgx"
)

func PrepareStatements() {
	database.DB.Prepare("getThreadBySlug", getThreadBySlug)
	database.DB.Prepare("getThreadById", getThreadById)

	database.DB.Prepare("getForumThreadsNsNd", getForumThreadsNsNd)
	database.DB.Prepare("getForumThreadsNsYd", getForumThreadsNsYd)
	database.DB.Prepare("getForumThreadsYsYd", getForumThreadsYsYd)
	database.DB.Prepare("getForumThreadsYsNd", getForumThreadsYsNd)
}

const getThreadBySlug = `
SELECT 
	id, 
	slug::TEXT, 
	author, 
	created, 
	forum::TEXT, 
	message, 
	title, 
	votes
FROM t_thread
WHERE slug=$1`

func (t *Thread) GetBySlug(slug interface{}) error {
	return database.DB.QueryRow("getThreadBySlug", slug).
		Scan(
			&t.Id,
			&t.Slug,
			&t.Author,
			&t.Created,
			&t.Forum,
			&t.Message,
			&t.Title,
			&t.Votes,
		)
}

const getThreadById = `
SELECT id, slug::TEXT, author, created, forum::TEXT, message, title, votes
FROM t_thread
WHERE id=$1::TEXT::INT`

func (t *Thread) GetById(id interface{}) error {
	return database.DB.QueryRow("getThreadById", id).
		Scan(
			&t.Id,
			&t.Slug,
			&t.Author,
			&t.Created,
			&t.Forum,
			&t.Message,
			&t.Title,
			&t.Votes,
		)
}

const createThread = `
INSERT INTO t_thread(slug, author, created, forum, message, title)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id`

func (t *Thread) Insert() error {
	tx := database.TxMustBegin()
	defer tx.Commit()

	err := tx.QueryRow(
		createThread,
		t.Slug,
		t.Author,
		t.Created,
		t.Forum,
		t.Message,
		t.Title,
	).Scan(&t.Id)

	if err != nil {
		tx.Rollback()
	}
	return err
}

const threadUpdate = `
UPDATE t_thread
SET message = COALESCE($1, message),
  title     = COALESCE($2, title)
WHERE id = $3
RETURNING message, title`

func (t *Thread) Update(tu ThreadUpdate) error {
	tx := database.TxMustBegin()
	defer tx.Commit()

	err := tx.QueryRow(threadUpdate, tu.Message, tu.Title, t.Id).
		Scan(&t.Message, &t.Title)
	if err != nil {
		tx.Rollback()
	}
	return err
}

const (
	getForumThreadsNsNd = `
SELECT id, slug::TEXT, author, created, forum::TEXT, message, title, votes
FROM t_thread
WHERE forum=$1
ORDER BY created
LIMIT $2::TEXT::BIGINT`

	getForumThreadsNsYd = `
SELECT id, slug::TEXT, author, created, forum::TEXT, message, title, votes
FROM t_thread
WHERE forum=$1
ORDER BY created DESC
LIMIT $2::TEXT::BIGINT`

	getForumThreadsYsNd = `
SELECT id, slug::TEXT, author, created, forum::TEXT, message, title, votes
FROM t_thread
WHERE forum=$1 AND created>=$2::TEXT::TIMESTAMPTZ
ORDER BY created
LIMIT $3::TEXT::BIGINT`

	getForumThreadsYsYd = `
SELECT id, slug::TEXT, author, created, forum::TEXT, message, title, votes
FROM t_thread
WHERE forum=$1 AND created<=$2::TEXT::TIMESTAMPTZ
ORDER BY created DESC
LIMIT $3::TEXT::BIGINT`
)

func (ts *ThreadSlice) Get(
	forumSlug interface{},
	limit []byte,
	since []byte,
	desc bool,
) error {
	var rows *pgx.Rows
	var err error

	if since == nil {
		if !desc {
			rows, err = database.DB.Query("getForumThreadsNsNd", forumSlug, limit)
		} else {
			rows, err = database.DB.Query("getForumThreadsNsYd", forumSlug, limit)
		}
	} else {
		if !desc {
			rows, err = database.DB.Query("getForumThreadsYsNd", forumSlug, since, limit)
		} else {
			rows, err = database.DB.Query("getForumThreadsYsYd", forumSlug, since, limit)
		}
	}

	if err != nil {
		return err
	}

	*ts = (*ts)[:100]
	l := 0
	for rows.Next() {
		err = rows.Scan(
			&(*ts)[l].Id,
			&(*ts)[l].Slug,
			&(*ts)[l].Author,
			&(*ts)[l].Created,
			&(*ts)[l].Forum,
			&(*ts)[l].Message,
			&(*ts)[l].Title,
			&(*ts)[l].Votes,
		)

		if err != nil {
			return err
		}
		l += 1
	}

	*ts = (*ts)[:l]

	return nil
}

const qByID = `
WITH find_diff AS (
	INSERT INTO t_vote(thread_id, nickname, vote) VALUES($1::TEXT::INT, (SELECT nickname FROM t_user WHERE nickname=$2), $3)
	ON CONFLICT ON CONSTRAINT t_vote_thread_id_nickname_idx DO
		UPDATE
		SET prev_vote=t_vote.vote,
			vote=excluded.vote
	RETURNING vote-prev_vote AS diff)
UPDATE t_thread
	SET votes=votes+(SELECT diff FROM find_diff)
	WHERE id=$1::TEXT::INT
RETURNING id, slug::TEXT, author, created, forum::TEXT, message, title, votes`

const qBySlug = `
WITH find_diff AS (
	INSERT INTO t_vote(thread_id, nickname, vote) VALUES((SELECT id FROM t_thread WHERE slug=$1), (SELECT nickname FROM t_user WHERE nickname=$2), $3)
	ON CONFLICT ON CONSTRAINT t_vote_thread_id_nickname_idx DO
		UPDATE
		SET prev_vote=t_vote.vote,
			vote=excluded.vote
	RETURNING vote-prev_vote AS diff)
UPDATE t_thread
	SET votes=votes+(SELECT diff FROM find_diff)
	WHERE slug=$1
RETURNING id, slug::TEXT, author, created, forum::TEXT, message, title, votes`

func (t *Thread) PutVote(slugOrId string, v Vote) error {
	tx := database.TxMustBegin()

	var row *pgx.Row
	if helpers.IsNumber(slugOrId) {
		row = tx.QueryRow(qByID, slugOrId, v.Nickname, v.Voice)
	} else {
		row = tx.QueryRow(qBySlug, slugOrId, v.Nickname, v.Voice)
	}

	err := row.Scan(&t.Id, &t.Slug, &t.Author, &t.Created, &t.Forum, &t.Message, &t.Title, &t.Votes)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

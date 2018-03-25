package thread

import (
	"github.com/VictoriaOtm/forum-api/database"
	"github.com/jackc/pgx"
	"bytes"
)

func (t *Thread) GetBySlug(slug string) error {
	return database.DBConnPool.QueryRow(database.GetThreadBySlug, slug).
		Scan(&t.Id, &t.Slug, &t.Author, &t.Created, &t.Forum, &t.Message, &t.Title, &t.Votes)
}

func (t *Thread) GetById(id string) error {
	return database.DBConnPool.QueryRow(database.GetThreadById, id).
		Scan(&t.Id, &t.Slug, &t.Author, &t.Created, &t.Forum, &t.Message, &t.Title, &t.Votes)
}

func (ta *Arr) Get(forumSlug string, limit []byte, since []byte, desc []byte) error {
	err := database.DBConnPool.QueryRow(database.GetRealForumSlug, forumSlug).Scan(&forumSlug)
	if err != nil {
		return database.ErrorForumNotExists
	}

	var rows *pgx.Rows
	d := bytes.Equal([]byte("true"), desc)

	if since == nil {
		if !d {
			rows, err = database.DBConnPool.Query(database.GetForumThreadsNsNd, forumSlug, limit)
		} else {
			rows, err = database.DBConnPool.Query(database.GetForumThreadsNsYd, forumSlug, limit)
		}
	} else {
		if !d {
			rows, err = database.DBConnPool.Query(database.GetForumThreadsYsNd, forumSlug, since, limit)
		} else {
			rows, err = database.DBConnPool.Query(database.GetForumThreadsYsYd, forumSlug, since, limit)
		}
	}

	if err != nil {
		return err
	}

	l := len(*ta)
	for rows.Next() {
		l += 1
		*ta = (*ta)[:l]

		err = rows.Scan(&(*ta)[l-1].Id, &(*ta)[l-1].Slug, &(*ta)[l-1].Author, &(*ta)[l-1].Created,
			&(*ta)[l-1].Forum, &(*ta)[l-1].Message, &(*ta)[l-1].Title, &(*ta)[l-1].Votes)

		if err != nil {
			return err
		}
	}

	return nil
}

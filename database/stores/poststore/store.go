package poststore

import (
	"bytes"

	"context"

	"errors"

	"log"

	"github.com/VictoriaOtm/forum-api/database"
	"github.com/jackc/pgx"
)

func PrepareStatements() {
	database.DB.Prepare("getPost", getPost)
	database.DB.Prepare("getPostThread", getPostThread)
	database.DB.Prepare("createPosts", createPosts)

	database.DB.Prepare("GetFlatNsNd", GetFlatNsNd)
	database.DB.Prepare("GetFlatNsYd", GetFlatNsYd)
	database.DB.Prepare("GetFlatYsNd", GetFlatYsNd)
	database.DB.Prepare("GetFlatYsYd", GetFlatYsYd)

	database.DB.Prepare("GetTreeNsNd", GetTreeNsNd)
	database.DB.Prepare("GetTreeNsYd", GetTreeNsYd)
	database.DB.Prepare("GetTreeYsNd", GetTreeYsNd)
	database.DB.Prepare("GetTreeYsYd", GetTreeYsYd)

	database.DB.Prepare("GetParentTreeNsNd", GetParentTreeNsNd)
	database.DB.Prepare("GetParentTreeNsYd", GetParentTreeNsYd)
	database.DB.Prepare("GetParentTreeYsNd", GetParentTreeYsNd)
	database.DB.Prepare("GetParentTreeYsYd", GetParentTreeYsYd)
}

const getPost = `
SELECT id, author, forum, thread, created, is_edited, parent, message
FROM t_posts
WHERE id=$1::TEXT::BIGINT`

func (p *Post) Get(id interface{}) error {
	return database.DB.QueryRow("getPost", id).
		Scan(
			&p.Id,
			&p.Author,
			&p.Forum,
			&p.Thread,
			&p.Created,
			&p.IsEdited,
			&p.Parent,
			&p.Message,
		)
}

const updatePost = `
UPDATE t_posts
SET is_edited = (SELECT (($1::TEXT IS NOT NULL OR is_edited) AND $1::TEXT != message)),
message = COALESCE($1, message)
WHERE id = $2
RETURNING message, is_edited`

func (p *Post) Update(pu PostUpdate) error {
	tx := database.TxMustBegin()
	defer tx.Commit()

	err := tx.QueryRow(updatePost, pu.Message, p.Id).
		Scan(&p.Message, &p.IsEdited)
	if err != nil {
		tx.Rollback()
	}

	return err
}

func (ps *PostSlice) Get(thrID int32, sort, since, limit []byte, desc bool) error {
	var rows *pgx.Rows
	var err error

	switch true {
	case bytes.Equal(sort, []byte("tree")):
		rows, err = getTreeSorted(thrID, limit, since, desc)
	case bytes.Equal(sort, []byte("parent_tree")):
		rows, err = getParentTreeSorted(thrID, limit, since, desc)
	case sort == nil, bytes.Equal(sort, []byte("flat")):
		rows, err = getFlatSorted(thrID, limit, since, desc)
	default:
		panic(sort)
	}

	if err != nil {
		panic(err)
	}

	l := 0
	*ps = (*ps)[:100]
	for rows.Next() {
		l++

		err = rows.Scan(
			&(*ps)[l-1].Id,
			&(*ps)[l-1].Author,
			&(*ps)[l-1].Forum,
			&(*ps)[l-1].Thread,
			&(*ps)[l-1].Created,
			&(*ps)[l-1].IsEdited,
			&(*ps)[l-1].Message,
			&(*ps)[l-1].Parent)
		if err != nil {
			return err
		}
	}
	*ps = (*ps)[:l]
	return nil
}

const GetFlatNsNd = `
SELECT id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
WHERE thread=$1
ORDER BY id
LIMIT $2::TEXT::INT`

const GetFlatNsYd = `
SELECT id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
WHERE thread=$1
ORDER BY id DESC
LIMIT $2::TEXT::INT`

const GetFlatYsNd = `
SELECT id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
WHERE thread=$1 AND id>$2::TEXT::BIGINT
ORDER BY id
LIMIT $3::TEXT::INT`

const GetFlatYsYd = `
SELECT id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
WHERE thread=$1 AND id<$2::TEXT::BIGINT
ORDER BY id DESC
LIMIT $3::TEXT::INT`

func getFlatSorted(thrID int32, limit, since []byte, desc bool) (*pgx.Rows, error) {
	if since == nil {
		if !desc {
			return database.DB.Query("GetFlatNsNd", thrID, limit)
		} else {
			return database.DB.Query("GetFlatNsYd", thrID, limit)
		}
	} else {
		if !desc {
			return database.DB.Query("GetFlatYsNd", thrID, since, limit)
		} else {
			return database.DB.Query("GetFlatYsYd", thrID, since, limit)
		}
	}
}

const (
	GetTreeNsNd = `
SELECT id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
WHERE thread=$1
ORDER BY parents
LIMIT $2::TEXT::INT`

	GetTreeNsYd = `
SELECT id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
WHERE thread=$1
ORDER BY parents DESC
LIMIT $2::TEXT::INT`

	GetTreeYsNd = `
SELECT id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
WHERE thread=$1 AND parents > (SELECT parents FROM t_posts WHERE id=$2::TEXT::BIGINT)
ORDER BY parents
LIMIT $3::TEXT::INT`

	GetTreeYsYd = `
SELECT id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
WHERE thread=$1 AND parents < (SELECT parents FROM t_posts WHERE id=$2::TEXT::BIGINT)
ORDER BY parents DESC
LIMIT $3::TEXT::INT`
)

func getTreeSorted(thrID int32, limit, since []byte, desc bool) (*pgx.Rows, error) {
	if since == nil {
		if !desc {
			return database.DB.Query("GetTreeNsNd", thrID, limit)
		} else {
			return database.DB.Query("GetTreeNsYd", thrID, limit)
		}
	} else {
		if !desc {
			return database.DB.Query("GetTreeYsNd", thrID, since, limit)
		} else {
			return database.DB.Query("GetTreeYsYd", thrID, since, limit)
		}
	}
}

const (
	GetParentTreeNsNd = `
SELECT t_posts.id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
JOIN (
	SELECT id
	FROM t_posts
	WHERE thread = $1 AND parent = 0
	ORDER BY id
	LIMIT $2::TEXT::INTEGER
) sub
ON sub.id = t_posts.main_parent
ORDER BY parents`

	GetParentTreeNsYd = `
SELECT t_posts.id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
JOIN (
	SELECT id
	FROM t_posts
	WHERE thread = $1 AND parent = 0
	ORDER BY id DESC
	LIMIT $2::TEXT::INTEGER
) sub
ON sub.id = t_posts.main_parent
ORDER BY t_posts.main_parent DESC, parents`

	GetParentTreeYsNd = `
SELECT t_posts.id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
JOIN (
	SELECT id
	FROM t_posts
	WHERE thread = $1 AND parent = 0 AND main_parent > (SELECT main_parent FROM t_posts WHERE id=$2::TEXT::BIGINT)
	ORDER BY id
	LIMIT $3::TEXT::INTEGER
) sub
ON sub.id = t_posts.main_parent
ORDER BY parents`

	GetParentTreeYsYd = `
SELECT t_posts.id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
JOIN (
	SELECT id
	FROM t_posts
	WHERE thread = $1 AND parent = 0 AND main_parent < (SELECT main_parent FROM t_posts WHERE id=$2::TEXT::BIGINT)
	ORDER BY id DESC
	LIMIT $3::TEXT::INTEGER
) sub
ON sub.id = t_posts.main_parent
ORDER BY t_posts.main_parent DESC, parents`
)

func getParentTreeSorted(thrID int32, limit, since []byte, desc bool) (*pgx.Rows, error) {
	if since == nil {
		if !desc {
			return database.DB.Query("GetParentTreeNsNd", thrID, limit)
		} else {
			return database.DB.Query("GetParentTreeNsYd", thrID, limit)
		}
	} else {
		if !desc {
			return database.DB.Query("GetParentTreeYsNd", thrID, since, limit)
		} else {
			return database.DB.Query("GetParentTreeYsYd", thrID, since, limit)
		}
	}
}

const getPostThread = `
SELECT thread, parents
FROM t_posts
WHERE id=$1`

var ErrorThreadsNotEqual = errors.New("threads not equal")

func (ps *PostSlice) ValidateParentsAndThread(threadID int32) error {
	batch := database.DB.BeginBatch()
	defer batch.Close()

	for _, post := range *ps {
		if post.Parent != 0 {
			batch.Queue("getPostThread", []interface{}{post.Parent}, nil, nil)
		}
	}

	err := batch.Send(context.Background(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	var selectedThreadID int32
	for i, post := range *ps {
		if post.Parent != 0 {
			var parents []int64
			err = batch.QueryRowResults().Scan(&selectedThreadID, &parents)
			if err != nil {
				return err
			}

			if selectedThreadID != threadID {
				return ErrorThreadsNotEqual
			}
			(*ps)[i].parents = parents
		}
	}

	return nil
}

const createPosts = `
INSERT INTO t_posts(id, author, forum, thread, created, parent, message, main_parent, parents) 
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`

const GetIds = `
SELECT array_agg(nextval('t_posts_id_seq')::BIGINT)
FROM generate_series(1,$1)`

func (ps *PostSlice) Insert() error {
	var ids []int64
	database.DB.QueryRow(GetIds, len(*ps)).Scan(&ids)

	batch := database.DB.BeginBatch()

	for i := range *ps {
		(*ps)[i].Id = ids[i]
		(*ps)[i].parents = append((*ps)[i].parents, (*ps)[i].Id)

		batch.Queue("createPosts", []interface{}{
			(*ps)[i].Id,
			(*ps)[i].Author,
			(*ps)[i].Forum,
			(*ps)[i].Thread,
			(*ps)[i].Created,
			(*ps)[i].Parent,
			(*ps)[i].Message,
			(*ps)[i].parents[0],
			(*ps)[i].parents,
		}, nil, nil)
	}

	batch.Send(context.Background(), nil)
	_, err := batch.ExecResults()
	if err != nil {
		log.Println(err)
	}

	batch.Close()
	return nil
}

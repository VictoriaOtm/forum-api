package post

import (
	"github.com/VictoriaOtm/forum-api/database"
	"github.com/VictoriaOtm/forum-api/models/user"
	"github.com/VictoriaOtm/forum-api/models/forum"
	"github.com/VictoriaOtm/forum-api/models/thread"
	"context"
	"github.com/jackc/pgx"
	"log"
)

func (p *PostDetails) Get(id string, relatedArr []string) error {
	err := database.DBConnPool.QueryRow("GetPost", id).
		Scan(&p.Post.Id, &p.Post.Author, &p.Post.Forum, &p.Post.Thread, &p.Post.Created,
		&p.Post.IsEdited, &p.Post.Parent, &p.Post.Message)
	if err != nil {
		return err
	}

	batch := database.DBConnPool.BeginBatch()
	for _, r := range relatedArr {
		switch r {
		case "user":
			batch.Queue("GetUserByNickname", []interface{}{p.Post.Author}, nil, nil)
		case "forum":
			batch.Queue("GetForumBySlug", []interface{}{p.Post.Forum}, nil, nil)
		case "thread":
			batch.Queue("GetThreadByIdInt", []interface{}{p.Post.Thread}, nil, nil)
		}
	}

	err = batch.Send(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, r := range relatedArr {
		switch r {
		case "user":
			p.User = &user.User{}
			err = batch.QueryRowResults().
				Scan(&p.User.Nickname, &p.User.Email, &p.User.Fullname, &p.User.About)
		case "forum":
			p.Forum = &forum.Forum{}
			err = batch.QueryRowResults().
				Scan(&p.Forum.Slug, &p.Forum.Title, &p.Forum.User, &p.Forum.Threads,
				&p.Forum.Posts)
		case "thread":
			p.Thread = &thread.Thread{}
			err = batch.QueryRowResults().
				Scan(&p.Thread.Id, &p.Thread.Slug, &p.Thread.Author, &p.Thread.Created,
				&p.Thread.Forum, &p.Thread.Message, &p.Thread.Title, &p.Thread.Votes)
		}
	}

	batch.Close()
	return nil
}

func (arr *PostsArr) Get(slugOrId, sort string, limit, since []byte, desc bool) error {
	thrID, err := getThreadID(slugOrId)
	if err != nil {
		return database.ErrorThreadNotExists
	}

	var rows *pgx.Rows
	switch sort {
	case "tree":
		rows, err = arr.getTreeSorted(thrID, limit, since, desc)
	case "parent_tree":
		rows, err = arr.getParentTreeSorted(thrID, limit, since, desc)
	default:
		rows, err = arr.getFlatSorted(thrID, limit, since, desc)
	}

	if err != nil {
		log.Fatalln(err)
	}

	arrLen := 0
	for rows.Next() {
		arrLen++
		*arr = (*arr)[:arrLen]

		err = rows.Scan(
			&(*arr)[arrLen-1].Id,
			&(*arr)[arrLen-1].Author,
			&(*arr)[arrLen-1].Forum,
			&(*arr)[arrLen-1].Thread,
			&(*arr)[arrLen-1].Created,
			&(*arr)[arrLen-1].IsEdited,
			&(*arr)[arrLen-1].Message,
			&(*arr)[arrLen-1].Parent)
		if err != nil {
			return err
		}
	}

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

func (arr *PostsArr) getFlatSorted(thrID int32, limit, since []byte, desc bool) (*pgx.Rows, error) {
	if since == nil {
		if !desc {
			return database.DBConnPool.Query(GetFlatNsNd, thrID, limit)
		} else {
			return database.DBConnPool.Query(GetFlatNsYd, thrID, limit)
		}
	} else {
		if !desc {
			return database.DBConnPool.Query(GetFlatYsNd, thrID, since, limit)
		} else {
			return database.DBConnPool.Query(GetFlatYsYd, thrID, since, limit)
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

func (arr *PostsArr) getTreeSorted(thrID int32, limit, since []byte, desc bool) (*pgx.Rows, error) {
	if since == nil {
		if !desc {
			return database.DBConnPool.Query(GetTreeNsNd, thrID, limit)
		} else {
			return database.DBConnPool.Query(GetTreeNsYd, thrID, limit)
		}
	} else {
		if !desc {
			return database.DBConnPool.Query(GetTreeYsNd, thrID, since, limit)
		} else {
			return database.DBConnPool.Query(GetTreeYsYd, thrID, since, limit)
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
ON sub.id = t_posts.parents[1]
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
ON sub.id = t_posts.parents[1]
ORDER BY t_posts.parents[1] DESC, parents`

	GetParentTreeYsNd = `
SELECT t_posts.id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
JOIN (
	SELECT id
	FROM t_posts
	WHERE thread = $1 AND parent = 0 AND parents[1] > (SELECT parents[1] FROM t_posts WHERE id=$2::TEXT::BIGINT)
	ORDER BY id
	LIMIT $3::TEXT::INTEGER
) sub
ON sub.id = t_posts.parents[1]
ORDER BY parents`

	GetParentTreeYsYd = `
SELECT t_posts.id, author, forum, thread, created, is_edited, message, parent
FROM t_posts
JOIN (
	SELECT id
	FROM t_posts
	WHERE thread = $1 AND parent = 0 AND parents[1] < (SELECT parents[1] FROM t_posts WHERE id=$2::TEXT::BIGINT)
	ORDER BY id DESC
	LIMIT $3::TEXT::INTEGER
) sub
ON sub.id = t_posts.parents[1]
ORDER BY t_posts.parents[1] DESC, parents`
)

func (arr *PostsArr) getParentTreeSorted(thrID int32, limit, since []byte, desc bool) (*pgx.Rows, error) {
	if since == nil {
		if !desc {
			return database.DBConnPool.Query(GetParentTreeNsNd, thrID, limit)
		} else {
			return database.DBConnPool.Query(GetParentTreeNsYd, thrID, limit)
		}
	} else {
		if !desc {
			return database.DBConnPool.Query(GetParentTreeYsNd, thrID, since, limit)
		} else {
			return database.DBConnPool.Query(GetParentTreeYsYd, thrID, since, limit)
		}
	}
}

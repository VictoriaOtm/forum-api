package database

const (
	GetThreadFromPost = `
SELECT thread, parents FROM t_posts WHERE id=$1`

	CreatePosts = `
INSERT INTO t_posts(id, author, forum, thread, created, parent, message, parents) 
VALUES($1, $2, $3, $4, $5, $6, $7, $8)`

	GetIds = `
SELECT array_agg(nextval('t_posts_id_seq')::BIGINT)
FROM generate_series(1,$1)`

	GetPost = `
SELECT id, author, forum, thread, created, is_edited, parent, message
FROM t_posts
WHERE id=$1::TEXT::BIGINT`
)

func postPrepareStatements() {
	mustPrepare("CreatePosts", CreatePosts)
	mustPrepare("GetThreadFromPost", GetThreadFromPost)
	mustPrepare("GetPost", GetPost)
}

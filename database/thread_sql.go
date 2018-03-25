package database

const (
	GetThreadBySlug = `
SELECT id, slug::TEXT, author, created, forum, message, title, votes
FROM t_thread
WHERE slug=$1`

	GetThreadById = `
SELECT id, slug::TEXT, author, created, forum, message, title, votes
FROM t_thread
WHERE id=$1::TEXT::INT`

	GetThreadByIdInt = `
SELECT id, slug::TEXT, author, created, forum, message, title, votes
FROM t_thread
WHERE id=$1`

	CreateThread = `
INSERT INTO t_thread(slug, author, created, forum, message, title)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id`

	GetForumThreadsNsNd = `
SELECT id, slug::TEXT, author, created, forum, message, title, votes
FROM t_thread
WHERE forum=$1
ORDER BY created
LIMIT $2::TEXT::BIGINT`

	GetForumThreadsNsYd = `
SELECT id, slug::TEXT, author, created, forum, message, title, votes
FROM t_thread
WHERE forum=$1
ORDER BY created DESC
LIMIT $2::TEXT::BIGINT`

	GetForumThreadsYsNd = `
SELECT id, slug::TEXT, author, created, forum, message, title, votes
FROM t_thread
WHERE forum=$1 AND created>=$2::TEXT::TIMESTAMPTZ
ORDER BY created
LIMIT $3::TEXT::BIGINT`

	GetForumThreadsYsYd = `
SELECT id, slug::TEXT, author, created, forum, message, title, votes
FROM t_thread
WHERE forum=$1 AND created<=$2::TEXT::TIMESTAMPTZ
ORDER BY created DESC
LIMIT $3::TEXT::BIGINT`

	GetThreadIDBySlug = `
SELECT id FROM t_thread WHERE slug=$1`

	GetThreadIDByID = `
SELECT id FROM t_thread WHERE id=$1::TEXT::INT`

	GetForumSlugByThreadID = `
SELECT forum FROM t_thread WHERE id=$1`
)

func prepareThreadStatements() {
	mustPrepare("GetThreadByIdInt", GetThreadByIdInt)
}

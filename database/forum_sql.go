package database

const (
	GetForumBySlug = `
SELECT slug::TEXT, title, f_user, threads, posts
FROM t_forum
WHERE slug=$1`

	GetRealForumSlug = `
SELECT slug::TEXT
FROM t_forum
WHERE slug=$1`

	CreateForum = `
INSERT INTO t_forum(slug, title, f_user)
VALUES($1, $2, $3)
`
	CreateForumUsers = `
INSERT INTO t_forum_user(slug, nickname, email, fullname, about)
VALUES($1, $2, $3, $4, $5)
ON CONFLICT DO NOTHING`

	UpdateForumPostsCount = `
UPDATE t_forum
SET posts = posts + $1
WHERE slug=$2`
)

func forumPrepareStaements() {
	mustPrepare("CreateForumUsers", CreateForumUsers)
	mustPrepare("GetForumBySlug", GetForumBySlug)
}

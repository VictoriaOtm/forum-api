package database

const (
	CreateUser = `
INSERT INTO t_user(nickname, email,fullname, about) 
VALUES ($1, $2, $3, $4) 
ON CONFLICT DO NOTHING`

	GetUserByNicknameOrEmail = `
SELECT nickname::TEXT, email::TEXT, fullname, about 
FROM t_user 
WHERE nickname=$1 OR email=$2`

	GetUserByNickname = `
SELECT nickname::TEXT, email::TEXT, fullname, about 
FROM t_user 
WHERE nickname=$1`

	GetUserByEmail = `
SELECT nickname::TEXT, email::TEXT, fullname, about
FROM t_user
WHERE email=$1`

	ReplaceUserInfo = `
UPDATE t_user 
SET email=COALESCE($1, email),
	fullname=COALESCE($2, fullname),
	about=COALESCE($3, about)
WHERE nickname=$4
RETURNING nickname::TEXT, email::TEXT, fullname, about`

	GetRealUserNickname = `
SELECT nickname::TEXT
FROM t_user
WHERE nickname=$1`

	GetForumUsersNsNd = `
SELECT nickname::TEXT, email::TEXT, fullname, about
FROM t_forum_user
WHERE slug=$1
ORDER BY nickname::citext
LIMIT $2::TEXT::INT`

	GetForumUsersNsYd = `
SELECT nickname::TEXT, email::TEXT, fullname, about
FROM t_forum_user
WHERE slug=$1
ORDER BY nickname::citext DESC
LIMIT $2::TEXT::INT`

	GetForumUsersYsNd = `
SELECT nickname::TEXT, email::TEXT, fullname, about
FROM t_forum_user
WHERE slug=$1 AND nickname>$2::TEXT::citext
ORDER BY nickname::citext
LIMIT $3::TEXT::INT`

	GetForumUsersYsYd = `
SELECT nickname::TEXT, email::TEXT, fullname, about
FROM t_forum_user
WHERE slug=$1 AND nickname<$2::TEXT::citext
ORDER BY nickname::citext DESC
LIMIT $3::TEXT::INT`
)

func userPrepareStatements() {
	mustPrepare("GetUserByNickname", GetUserByNickname)
}

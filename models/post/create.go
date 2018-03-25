package post

import (
	"github.com/VictoriaOtm/forum-api/helpers"
	"github.com/VictoriaOtm/forum-api/database"
	"context"
	"github.com/VictoriaOtm/forum-api/models/user"
	"strings"
	"time"
)

func getThreadID(slugOrId string) (id int32, err error) {
	if helpers.IsNumber(slugOrId) {
		err = database.DBConnPool.QueryRow(database.GetThreadIDByID, slugOrId).Scan(&id)
	} else {
		err = database.DBConnPool.QueryRow(database.GetThreadIDBySlug, slugOrId).Scan(&id)
	}

	return
}

func (arr *PostsArr) Create(slugOrId string) error {
	tx, err := database.DBConnPool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	threadID, err := getThreadID(slugOrId)
	if err != nil {
		return database.ErrorThreadNotExists
	}

	if len(*arr) == 0 {
		return nil
	}

	var forumSlug string
	err = tx.QueryRow(database.GetForumSlugByThreadID, threadID).Scan(&forumSlug)
	if err != nil {
		return database.ErrorForumNotExists
	}

	usersMap := make(map[string]string)
	for _, post := range *arr {
		usersMap[strings.ToLower(post.Author)] = post.Author
	}

	batch := tx.BeginBatch()
	defer batch.Close()

	usersCount := 0
	for nickname := range usersMap {
		usersCount++
		batch.Queue("GetUserByNickname", []interface{}{nickname}, nil, nil)
	}

	for _, post := range *arr {
		if post.Parent != 0 {
			batch.Queue("GetThreadFromPost", []interface{}{post.Parent}, nil, nil)
		}
	}
	batch.Send(context.Background(), nil)

	userArr := make(user.Arr, usersCount)
	for i := 0; i < usersCount; i++ {
		err = batch.QueryRowResults().Scan(&userArr[i].Nickname, &userArr[i].Email, &userArr[i].Fullname, &userArr[i].About)
		if err != nil {
			return database.ErrorUserNotExists
		}
		usersMap[strings.ToLower(userArr[i].Nickname)] = userArr[i].Nickname
	}

	var parentThread int32
	parentsArr := make([][]int64, len(*arr))
	for i, post := range *arr {
		if post.Parent == 0 {
			continue
		}

		batch.QueryRowResults().Scan(&parentThread, &parentsArr[i])
		if parentThread != threadID {
			return database.ErrorPostParentConflict
		}
	}

	ids := make([]int64, len(*arr))
	err = tx.QueryRow(database.GetIds, len(*arr)).Scan(&ids)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	for i, post := range *arr {
		(*arr)[i].Id = ids[i]
		(*arr)[i].Forum = forumSlug
		(*arr)[i].Thread = threadID
		(*arr)[i].Created = currentTime
		(*arr)[i].Author = usersMap[strings.ToLower(post.Author)]
		parentsArr[i] = append(parentsArr[i], ids[i])

		batch.Queue("CreatePosts", []interface{}{(*arr)[i].Id, (*arr)[i].Author, (*arr)[i].Forum, (*arr)[i].Thread, (*arr)[i].Created, (*arr)[i].Parent, (*arr)[i].Message, parentsArr[i]}, nil, nil)
	}

	for _, usr := range userArr {
		batch.Queue("CreateForumUsers", []interface{}{forumSlug, usr.Nickname, usr.Email, usr.Fullname, usr.About}, nil, nil)
	}

	batch.Send(context.Background(), nil)

	for range *arr {
		_, err = batch.ExecResults()
		if err != nil {
			return err
		}
	}

	for range userArr {
		_, err = batch.ExecResults()
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(database.UpdateForumPostsCount, len(*arr), forumSlug)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

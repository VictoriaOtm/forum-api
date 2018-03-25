package database

import "fmt"

var (
	ErrorForumNotExists = fmt.Errorf("ErrorForumNotExists")
	ErrorForumConflict  = fmt.Errorf("ErrorForumConflict")

	ErrorUserConflict = fmt.Errorf("user exists")
	ErrorUserNotExists = fmt.Errorf("ErrorUserNotExists")

	ErrorThreadConflict = fmt.Errorf("ErrorThreadConflict")
	ErrorThreadNotExists = fmt.Errorf("ErrorThreadNotExists")

	ErrorPostParentConflict = fmt.Errorf("ErrorPostParentConflict")
)

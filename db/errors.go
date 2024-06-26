package db

import "fmt"

type databaseConnectionError struct {
	Err error
}

func (e *databaseConnectionError) Error() string {
	return fmt.Sprintf("Error connecting the database: \n%s", e.Err.Error())
}

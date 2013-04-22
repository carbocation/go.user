/*
To be initialized, this class expects a live Postgres database connection handle to 
be passed in. For example:

var db *sql.DB

func main() {
	db = (...get the object...)

	user.Initialize(db)
}
*/
package user

import (
	"database/sql"
)

//Define package-level globals
var db *sql.DB

//This function must be called in the main package implementing the forum
func Initialize(pool *sql.DB) {
	db = pool
}

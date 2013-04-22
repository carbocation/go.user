package user

var queries = struct {
	UserFindById string
	UserFindByHandle	string
	UserCreate string
}{
	UserFindById: `SELECT id, handle, email, password, created FROM account WHERE id=$1 LIMIT 1`,
	UserFindByHandle: `SELECT id, handle, email, password, created FROM account WHERE handle=$1 LIMIT 1`,
	UserCreate: `INSERT INTO account (handle, email, password) VALUES ($1, $2, $3) RETURNING id`, 
}

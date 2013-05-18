package user

import (
	"errors"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
)

type User struct {
	Id       int64     "The user's auto-incremented ID"
	Handle   string    "The user's name"
	Email    string    "The user's email"
	Password string    `schema:"-"` //Do not allow this to be set by users
	Created  time.Time "The creation timestamp of the user's account"

	//Values below here should not be stored in the database
	PlaintextPassword string "Plaintext password. Not stored in the database; used for logging in."
}

func (u *User) GetId() int64 {
	return u.Id
}

func (u *User) Guest() bool {
	if u.Id == 0 {
		return true
	}

	return false
}

func (u *User) Register() (err error) {
	//Make sure they sent us a bcrypt-able password
	err = u.SetPassword(u.PlaintextPassword)
	if err != nil {
		//Fail here to avoid modifying the DB
		return
	}

	//We will now try to store this new user to the database.
	err = u.createInDb()

	return
}

func (u *User) createInDb() error {
	//Wrap in a transaction
	tx, err := db.Begin()

	CreateUserStmt, err := tx.Prepare(queries.UserCreate)
	if err != nil {
		tx.Rollback()
		return errors.New("Error: We had a database problem.")
	}
	defer CreateUserStmt.Close()

	//Note: because pq handles LastInsertId oddly (or not at all?), instead of 
	//calling .Exec() then .LastInsertId, we prepare a statement that ends in 
	//`RETURNING id` and we .QueryRow().Select() the result  
	err = CreateUserStmt.QueryRow(u.Handle, u.Email, u.Password).Scan(&u.Id)
	if err != nil {
		tx.Rollback()
		return errors.New("Error: your username or email address was already found in the database. Please choose differently.")
	}

	//Declare transactional victory
	tx.Commit()

	return nil
}

//SetPassword takes a plaintext password and hashes it with bcrypt and sets the
//password field to the hash.
func (u *User) SetPassword(password string) error {
	if len(password) < 1 {
		return errors.New("Error: no password was provided.")
	}

	hpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	//Regardless of success or failure, we now get rid of the PlaintextPassword value
	u.PlaintextPassword = ""

	if err != nil {
		return err
	}
	u.Password = string(hpass)

	return nil
}

//If a login checks out, return a new user object
func (login *User) Login() (user *User, err error) {
	user, err = FindOneByHandle(login.Handle)
	//If a valid user couldn't be found, the user shall become logged out
	if err != nil {
		user = new(User)
		err = errors.New("The username or password was invalid")
		return
	}

	//If the user account's password validates, update u with the new user object
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.PlaintextPassword))

	//Regardless of success or failure, we now get rid of the PlaintextPassword value
	login.PlaintextPassword = ""

	if err != nil {
		user = new(User)
		err = errors.New("The username or password was invalid")
		return
	}

	return
}

func FindOneByHandle(handle string) (user *User, err error) {
	//Initialize empty user
	user = new(User)

	//Find one or zero existent users
	FindUserStmt, err := db.Prepare(queries.UserFindByHandle)
	if err != nil {
		return
	}
	defer FindUserStmt.Close()

	//Read any found values into the user object
	err = FindUserStmt.QueryRow(handle).Scan(&user.Id, &user.Handle, &user.Email, &user.Password, &user.Created)
	if err != nil {
		user = new(User)
	}

	//If there are no records, the error will be non-nil
	//You may want to obscure this error message elsewhere, 
	//E.g., to hide whether a login failure was due to a non-existent
	//user or a bad password for a real user.  
	return
}

//If a user with this id exists, return a *User object. Otherwise, nil.
func FindOneUserById(id int64) (user *User, err error) {
	//Initialize empty user
	user = new(User)

	//Find one or zero existent users
	FindUserStmt, err := db.Prepare(queries.UserFindById)
	if err != nil {
		return
	}
	defer FindUserStmt.Close()

	//Read any found values into the user object
	err = FindUserStmt.QueryRow(id).Scan(&user.Id, &user.Handle, &user.Email, &user.Password, &user.Created)
	if err != nil {
		user = new(User)
	}

	//If there are no records, the error will be non-nil
	//You may want to obscure this error message elsewhere, 
	//E.g., to hide whether a login failure was due to a non-existent
	//user or a bad password for a real user.  
	return
}

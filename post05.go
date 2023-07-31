package post05

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

var (
	Hostname = ""
	Port = 2345
	Username = ""
	Password = ""
	Database = ""
)

type Userdata struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

func openConnection() (*sql.DB, error) {
	// connection string
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)

	// open database
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// The function returns the User ID of the username
// -1 if the user does not exist
func exist(username string) int  {
	username = strings.ToLower(username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := -1
	statement := fmt.Sprintf(`SELECT "id" FROM "users" where username = '%s'`, username)
	rows, err := db.Query(statement)

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
			return -1
		}
		userID = id
	}
	defer rows.Close()
	return userID
}


// AddUser adds a new user to the database
// Returns new User ID
// -1 if there wa an error
func AddUser(d Userdata) int  {
	d.Username = strings.ToLower(d.Username)
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := exist(d.Username)
	if userID != -1 {
		fmt.Println("User already exists:", Username)
		return -1
	}

	insertStatement := `insert into "users" ("username") values ($1)`
	_, err = db.Exec(insertStatement, d.Username)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	userID = exist(d.Username)
	if userID == -1 {
		return userID
	}

	insertStatement = `insert into "userdata" ("userid", "name", "surname", "description") values ($1, $2, $3, $4)`

	_, err = db.Exec(insertStatement, userID, d.Name, d.Surname, d.Description)
	if err != nil {
		fmt.Println("db.Exec():", err)
		return -1
	}

	return userID
}

// DeleteUser deletes an existing user
func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	statement := fmt.Sprintf(`SELECT "username" FROM "users" WHERE id=%d`, id)
	rows, err := db.Query(statement)

	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}
	defer rows.Close()

	if exist(username) != id {
		return fmt.Errorf("User with ID %d does not exist", id)
	}

	deleteStatement := `DELETE FROM "userdata" WHERE userid=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}

	deleteStatement = `DELETE FROM "users" where id=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}

	return nil
}

//
func ListUsers() ([]Userdata, error) {
	Data := []Userdata{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query(`select "id", "username", "name", "surname", "description" FROM "users", "userdata" WHERE users.id = userdata.userid`)
	if err != nil {
		return Data, err
	}

	for rows.Next() {
		var id int
		var username string
		var name string
		var surname string
		var description string
		err = rows.Scan(&id, &username, &name, &surname, &description)
		temp := Userdata{ID: id, Username: username, Name: name, Surname: surname, Description: description}

		Data = append(Data, temp)
		if err != nil {
			return Data, err
		}
	}
	defer rows.Close()
	return Data, nil
}

// UpdateUser is for updating an existing user
func UpdateUser(d Userdata) error  {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	userID := exist(d.Username)
	if userID == -1 {
		return errors.New("User does not exist")
	}
	d.ID = userID
	updateStatement := `UPDATE "userdata" SET "name"=$1, "surname"=$2, "description"=$3 WHERE "userid"=$4`
	_, err = db.Exec(updateStatement, d.Name, d.Surname, d.Description, d.ID)
	if err != nil {
		return err
	}

	return nil
}
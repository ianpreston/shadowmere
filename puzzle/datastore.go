package puzzle

import (
	_ "github.com/lib/pq"
	"database/sql"
	"code.google.com/p/go.crypto/bcrypt"
)

type Datastore struct {
	db *sql.DB
}

type RegisteredNick struct {
	Id int
	Nick string
	Email string
	Passwd string
}

func NewDatastore(pgUrl string) (*Datastore, error) {
	db, err := sql.Open("postgres", pgUrl)
	if err != nil {
		return nil, err
	}

	ds := &Datastore{
		db: db,
	}
	return ds, nil
}

func (ds *Datastore) GetRegisteredNick(nick string) (*RegisteredNick, error) {
	rn := &RegisteredNick{}

	row := ds.db.QueryRow(
		"SELECT id, nick, email, passwd FROM RegisteredNicks WHERE nick = $1",
		nick,
	)
	err := row.Scan(&rn.Id, &rn.Nick, &rn.Email, &rn.Passwd)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	return rn, nil
}

func (ds *Datastore) Authenticate(rn *RegisteredNick, passwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(rn.Passwd), []byte(passwd))
	if err == nil {
		return true
	} else {
		return false
	}
}

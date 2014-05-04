package shadowmere

import (
	_ "github.com/lib/pq"
	"database/sql"
	"code.google.com/p/go.crypto/bcrypt"
)

/**
 * Persistent services state - which nicks and channels are registered,
 * access control list settings, and so on.
 */
type Datastore struct {
	db *sql.DB

	RegisteredNicks *RegisteredNickRepo
}

type RegisteredNickRepo struct {
	ds *Datastore
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
	ds.RegisteredNicks = &RegisteredNickRepo{ds}

	return ds, nil
}

func (repo *RegisteredNickRepo) GetByNick(nick string) (*RegisteredNick, error) {
	rn := &RegisteredNick{}

	row := repo.ds.db.QueryRow(
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

func (repo *RegisteredNickRepo) Authenticate(rn *RegisteredNick, passwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(rn.Passwd), []byte(passwd))
	if err == nil {
		return true
	} else {
		return false
	}
}

func (repo *RegisteredNickRepo) Register(rn *RegisteredNick) error {
	passwd, err := bcrypt.GenerateFromPassword([]byte(rn.Passwd), 12)
	if err != nil {
		return err
	}

	var registeredNickId int
	err = repo.ds.db.QueryRow(
		"INSERT INTO RegisteredNicks (nick, email, passwd) VALUES ($1, $2, $3) RETURNING id;",
		rn.Nick,
		rn.Email,
		passwd,
	).Scan(&registeredNickId)
	if err != nil {
		return err
	}

	return nil
}

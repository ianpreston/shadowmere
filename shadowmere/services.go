package shadowmere


type Services struct {
	connection *Connection

	nickserv *NickServ
	chanserv *ChanServ

	datastore *Datastore
	curstate *CurState
}

func NewServices(pgUrl, name, addr, pass string) (*Services, error) {
	mere := &Services{}

	datastore, err := NewDatastore(pgUrl)
	if err != nil {
		return nil, err
	}

	connection, err := NewConnection(mere, name, addr, pass)
	if err != nil {
		return nil, err
	}

	mere.datastore = datastore
	mere.connection = connection
	mere.curstate = NewCurState()
	mere.nickserv = NewNickserv(mere)
	mere.chanserv = NewChanserv(mere)

	return mere, nil
}

func (mere *Services) Start() {
	mere.connection.Start()
}
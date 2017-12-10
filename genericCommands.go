package ts3Query

// Help sends to the writer the help command with the command supplied added to the end
func (t *Ts3Query) Help(command string) (response string, err error) {
	_, err = t.rw.Write([]byte("help " + escapeString(command)))
	if err != nil {
		return "", err
	}
	response, err = t.readResponse()
	return
}

// Use allows the query to direct commands to a specific server,
// it must be done before any specific server commands can be done
func (t *Ts3Query) Use(virtualServerID string) error {
	_, err := t.rw.Write([]byte("use " + escapeString(virtualServerID)))
	if err != nil {
		return err
	}
	res, err := t.readResponse()
	if err != nil {
		return err
	}

	return getError(res)
}

// Login takes a username and a password for the teamspeak server. It format's the input and writes it to the writer
func (t *Ts3Query) Login(username, password string) error {
	_, err := t.rw.Write([]byte("login " + escapeString(username) + " " + escapeString(password) + "\n"))
	if err != nil {
		return err
	}
	s, err := t.readResponse()
	if err != nil {
		return err
	}
	return getError(s)

}

// Logout writes the logout command to the writer
func (t *Ts3Query) Logout() error {
	_, err := t.rw.Write([]byte("logout"))
	if err != nil {
		return err
	}
	s, err := t.readResponse()
	return getError(s)
}

// Quit will write the command to quit the session
func (t *Ts3Query) Quit() error {
	_, err := t.rw.Write([]byte("quit"))
	if err != nil {
		return err
	}
	s, err := t.readResponse()
	return getError(s)
}

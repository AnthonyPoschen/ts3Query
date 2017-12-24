package ts3Query

import (
	"fmt"
	"strings"
)

type ServerGroup struct {
	SGID string
	Name string
}

//ServerGroupList returns the raw string of the server group list
func (t *Ts3Query) ServerGroupList() (ServerGroups []ServerGroup, err error) {
	err = t.sendMessage("servergrouplist")
	if err != nil {
		return
	}
	res, err := t.readResponse()
	if err != nil {
		return
	}
	groups := strings.Split(res, "|")

	for _, v := range groups {
		items := strings.Split(v, " ")
		m := make(map[string]string)
		for _, item := range items {
			i := strings.SplitN(item, "=", 2)
			if len(i) != 2 {
				err = fmt.Errorf("Unable to get key value pair for: %s", item)
				return
			}
			//fmt.Printf("1: %s\n2: %s\n", i[0], i[1])
			m[i[0]] = i[1]
		}
		ServerGroups = append(ServerGroups, ServerGroup{Name: unEscapeString(m["name"]), SGID: m["sgid"]})
	}
	return
}

func (t *Ts3Query) ServerGroupClientList(groupID string) (userIDs []string, err error) {
	err = t.sendMessage("servergroupclientlist sgid=" + groupID)
	if err != nil {
		return
	}
	res, err := t.readResponse()
	if err != nil {
		return
	}

	err = getError(res)
	if err != nil {
		return
	}
	results := strings.Split(res, "\n")
	results = strings.Split(results[0], "|")
	for _, item := range results {
		if strings.Contains(item, "error id=") {
			continue
		}
		pair := strings.SplitN(item, "=", 2)
		if len(pair) == 2 {
			if pair[0] == "cldbid" {
				userIDs = append(userIDs, pair[1])
			}
		}
	}
	return
}

// ServerGroupAddClient Adds a user to a server group
func (t *Ts3Query) ServerGroupAddClient(clientID string, groupID string) error {
	err := t.sendMessage(fmt.Sprintf("servergroupaddclient sgid=%s cldbid=%s", groupID, clientID))
	if err != nil {
		return err
	}
	res, err := t.readResponse()
	if err != nil {
		return err
	}
	return getError(res)
}

// ServerGroupDelClient removes a user from a server group
func (t *Ts3Query) ServerGroupDelClient(clientID string, groupID string) error {
	err := t.sendMessage(fmt.Sprintf("servergroupdelclient sgid=%s cldbid=%s", groupID, clientID))
	if err != nil {
		return err
	}
	res, err := t.readResponse()
	if err != nil {
		return err
	}
	return getError(res)
}

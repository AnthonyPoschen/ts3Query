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
		fmt.Println("Name:", m["name"], "||", unEscapeString(m["name"]))
		ServerGroups = append(ServerGroups, ServerGroup{Name: unEscapeString(m["name"]), SGID: m["sgid"]})
	}

	fmt.Printf("%+v", ServerGroups)
	return
}

package ts3Query

import (
	"fmt"
	"strconv"
	"strings"
)

// Client is the ts3 DB Client
type Client struct {
	DBID          string
	UUID          string
	Name          string
	Created       string
	LastConnected string
}

//ClientDBList returns a list of all the clients that have been connected to the teamspeak server
func (t *Ts3Query) ClientDBList() []Client {
	first := true
	var ClientList []Client
	var lastID string
	var NumberOfClients int
	for {
		msg := "clientdblist"
		if first == false {
			fmt.Println("Doing scan from", lastID)
			msg += " start=" + lastID
		}
		first = false
		t.sendMessage(msg)

		res, err := t.readResponse()
		if err != nil {
			fmt.Printf("error:%s\n", err)
			break
		}

		err = getError(res)
		if err != nil {
			break
		}
		groups := strings.Split(res, "|")
		for _, v := range groups {
			items := strings.Split(v, " ")
			m := make(map[string]string)
			for _, item := range items {
				i := strings.SplitN(item, "=", 2)
				if len(i) <= 1 {
					break
				}
				//fmt.Printf("1: %s\n2: %s\n", i[0], i[1])
				m[i[0]] = i[1]
			}
			if _, ok := m["cldbid"]; ok {
				ClientList = append(ClientList, Client{DBID: m["cldbid"], UUID: m["client_unique_identifier"], Name: unEscapeString(m["client_nickname"]), Created: m["client_created"], LastConnected: m["client_lastconnected"]})
			}

		}
		if NumberOfClients == len(ClientList) {
			break
		}
		NumberOfClients = len(ClientList)
		// it goes by count as people can be removed and not by the ID.
		lastID = strconv.Itoa(len(ClientList))

	}

	return ClientList
}

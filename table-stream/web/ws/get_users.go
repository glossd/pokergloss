package ws

import "github.com/glossd/pokergloss/auth/authid"

func GetTableUsers(tableID string) []*authid.Identity {
	load, ok := tableConns.Load(tableID)
	if ok {
		conns := load.(*Hub).clients
		users := make([]*authid.Identity, 0, len(conns))
		for client := range conns {
			users = append(users, client.iden)
		}
		return users
	}
	return nil
}

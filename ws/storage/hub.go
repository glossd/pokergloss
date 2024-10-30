package storage

type Hub interface {
	Register(conn Conn)
	Unregister(conn Conn)
}

package store

type Store interface {
	Writer() (Writer, error)
	Reader() (Reader, error)
	//Watcher() (Watcher, error)
	Close() error
}

type Writer interface {
	//Write(string, map[string]interface{}) error
	PutNode(string, string) error
	DeleteNode(string, string) error
	Close() error
}

type Reader interface {
	//Read(string) (map[string]interface{}, error)
	GetNode(string, string) (map[string]interface{}, error)
	Close() error
}

//type Watcher interface {
//	Watch()
//}

package persist

//Db interface enable switching between different persistent storage
type Db interface {
	SetNX(string, interface{}, int64) error
	Incr(string) (int64, error)
	Reset(string, int) error
}

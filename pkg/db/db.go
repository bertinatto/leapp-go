package db

var (
	DB = make(map[uint32]string)
)

func Set(k uint32, v string) {
	DB[k] = v
}
func Get(k uint32) string {
	return DB[k]
}

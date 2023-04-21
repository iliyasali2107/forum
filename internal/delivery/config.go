package delivery

type config struct {
	port int
	env  string
	db   struct {
		connectionString string
	}
}

package main

func main() {
	conf := ParseConfig()
	pool := NewDB(conf)
	Migrate(pool)
}

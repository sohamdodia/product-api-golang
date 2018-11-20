package config

type constants struct {
	PORT  int
	Mongo mongo
}

type mongo struct {
	URL    string
	DBName string
}

var Constants = constants{
	PORT: 8080,
	Mongo: mongo{
		URL:    "localhost",
		DBName: "products",
	},
}

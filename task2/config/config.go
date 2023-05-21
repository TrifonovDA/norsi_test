package config

type BdCredentials struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

var BdCred = BdCredentials{
	Host:     "localhost",
	Port:     "5432",
	Database: "norsi",
	Username: "d.triphonov",
	Password: "-------",
}

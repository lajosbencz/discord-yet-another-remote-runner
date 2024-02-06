package discordyetanotherremoterunner

type ConfigServer struct {
	Name  string
	Start string
	Stop  string
}

type Config struct {
	Guild   string
	Servers []ConfigServer
}

package server

type Server struct {
	Ip   string `yaml:"ip"`
	User string `yaml:"user"`
	Pass string `yaml:"password"`
}

type TagetServers struct {
	Servers []Server `yaml:"target_servers"`
}

func NewServer(ip, user, password string) Server {
	return Server{
		Ip:   ip,
		User: user,
		Pass: password,
	}
}

func NewTargetServers() TagetServers {
	return TagetServers{
		Servers: make([]Server, 0),
	}
}

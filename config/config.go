package config

type Config struct {
	Domains []struct {
		Name string
		Ns   []string
	}
}

type domain struct {
	Name string
	Ns   []string
}

var OreOreConfig = Config{
	Domains: []struct {
		Name string
		Ns   []string
	}{
		{
			Name: "oreore.net.",
			Ns: []string{
				"ns-973.awsdns-57.net.",
				"ns-409.awsdns-51.com.",
				"ns-1025.awsdns-00.org.",
				"ns-1854.awsdns-39.co.uk.",
			},
		},
		{
			Name: "oreore.dev.",
			Ns:   []string{},
		},
		{
			Name: "oreore.app.",
			Ns:   []string{},
		},
		{
			Name: "oreore.page.",
			Ns:   []string{},
		},
	},
}

package gorealdebrid

type RealDebrid struct {
	Client *RealDebridClient
}

const baseUrl = "https://api.real-debrid.com/rest/1.0/"

func NewRealDebrid(apiKey string) *RealDebrid {
	return &RealDebrid{
		Client: &RealDebridClient{
			ApiKey: apiKey,
		},
	}
}

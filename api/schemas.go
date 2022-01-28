package api

type route struct {
	name   string
	method map[string]Handler
}

var routeSchema = []route{
	{
		name: "/ping",
		method: map[string]Handler{
			"GET":  getPingHandler,
			"POST": postPingHandler,
		},
	},
	{
		name: "/stats",
		method: map[string]Handler{
			"GET": getStatsHandler,
		},
	},
	{
		name: "/detail",
		method: map[string]Handler{
			"GET": getDetailHandler,
		},
	},
}

type postPingRequest struct {
	Name string `json:"service_name"`
	URL  string `json:"service_url"`
}

type getDetailRequest postPingRequest

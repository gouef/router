package router

type RouteList struct {
	list []struct {
		Router Router
		Rank   int
	}

	ranks  [][]Router
	domain *string
	path   *string
}

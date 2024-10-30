package paging

type Page struct {
	Skip  int64
	Limit int64
}

func DefaultPage() Page {
	return Page{
		Skip:  0,
		Limit: 20,
	}
}

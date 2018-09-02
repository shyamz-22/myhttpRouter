package router

type PathParams []param

func (params PathParams) ByName(name string) string {
	for _, p := range params {
		if p.key == name {
			return p.value
		}
	}

	return ""
}

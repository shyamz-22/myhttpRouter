package router

type Param struct {
	Key   string
	Value string
}

type PathParams []Param

func (params PathParams) ByName(name string) string {
	for _, p := range params {
		if p.Key == name {
			return p.Value
		}
	}

	return ""
}

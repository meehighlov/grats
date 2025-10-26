package inlinekeyboard

type Builder struct {
	keyboard []*[]map[string]interface{}
}

func New() *Builder {
	return &Builder{}
}

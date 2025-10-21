package ports

type WorkFlow interface {
	Save(k string, v map[string]any) error
	Load() ([]map[string]any, error)
}
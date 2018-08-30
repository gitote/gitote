package git

func (repo *Repository) GetHook(name string) (*Hook, error) {
	return GetHook(repo.Path, name)
}

func (repo *Repository) Hooks() ([]*Hook, error) {
	return ListHooks(repo.Path)
}

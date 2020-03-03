package vm

type PaginatedProjects struct {
	Pagination  *Pagination `json:"pagination"`
	ProjectList []*Project  `json:"project_list"`
}

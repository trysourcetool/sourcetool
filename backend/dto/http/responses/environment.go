package responses

type EnvironmentResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Color     string `json:"color"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type ListEnvironmentsResponse struct {
	Environments []*EnvironmentResponse `json:"environments"`
}

type GetEnvironmentResponse struct {
	Environment *EnvironmentResponse `json:"environment"`
}

type CreateEnvironmentResponse struct {
	Environment *EnvironmentResponse `json:"environment"`
}

type UpdateEnvironmentResponse struct {
	Environment *EnvironmentResponse `json:"environment"`
}

type DeleteEnvironmentResponse struct {
	Environment *EnvironmentResponse `json:"environment"`
}

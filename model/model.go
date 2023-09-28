package model

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
}

type Owner struct {
	Login string `json:"login"`
}

type Repositories []Repository

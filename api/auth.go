package api

type authOrgInfo struct {
	Key     string            `json:"_key"`
	Name    string            `json:"name"`
	Members map[string]string `json:"members"`
}

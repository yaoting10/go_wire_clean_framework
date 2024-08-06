package response

type VersionResp struct {
	Version string `json:"version"`
	Url     string `json:"url"`
	Content string `json:"content"`
	Force   bool   `json:"force"`
}

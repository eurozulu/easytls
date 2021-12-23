package templates

type Page struct {
	PageTitle    string    `json:"title"`
	SubTitle     string    `json:"sub_title,omitempty"`
	Paragraphs   []string  `json:"paragraphs,omitempty"`
	Options      []*Option `json:"options,omitempty"`
	TemplateName string    `json:"templateName"`
}

type Option struct {
	Title  string `json:"title"`
	Target string `json:"target"`
}

package sp

type Result struct {
	ActionResult actionResult `xml:"action_result"`
	Reference    string       `xml:"reference"`
	RequestID    string       `xml:"request_id"`
}

type actionResult struct {
	Status      int         `xml:"status"`
	Code        int         `xml:"code"`
	Detail      string      `xml:"detail"`
	RedirectURL redirectURL `xml:"redirect"`
}
type redirectURL struct {
	URL string `xml:"url"`
}

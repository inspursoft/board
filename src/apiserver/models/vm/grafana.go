package vm

type evalMatch struct {
	Metric string                 `json:"metric"`
	Tags   map[string]interface{} `json:"tags"`
	Value  interface{}            `json:"value"`
}
type GrafanaNotification struct {
	Title       string      `json:"title"`
	RuleID      int         `json:"ruleId"`
	RuleName    string      `json:"ruleName"`
	RuleURL     string      `json:"ruleUrl"`
	State       string      `json:"state"`
	ImageURL    string      `json:"imageUrl"`
	Message     string      `json:"message"`
	EvalMatches []evalMatch `json:"evalMatches"`
}

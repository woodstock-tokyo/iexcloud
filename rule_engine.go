package iex

type (
	Rule struct {
		Token      string      `json:"token"`
		RuleSet    string      `json:"ruleSet"`
		Type       string      `json:"type"`
		RuleName   string      `json:"ruleName"`
		Conditions []Condition `json:"conditions"`
		Outputs    []Output    `json:"outputs"`
	}

	// Output webhook case
	Output struct {
		Frequency int    `json:"frequency"`
		Method    string `json:"method"`
		Url       string `json:"url"`
	}
	Condition []any

	CreatedRuleResponse struct {
		ID     string `json:"id"`
		Weight uint   `json:"weight"`
	}
)

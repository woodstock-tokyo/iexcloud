package iex

type (
	Rule struct {
		Token      string      `json:"token"`
		ID         string      `json:"id,omitempty"`
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

	RequestRule struct {
		Token  string `json:"token"`
		RuleID string `json:"ruleId"`
	}

	CreatedRuleResponse struct {
		ID     string `json:"id"`
		Weight uint   `json:"weight"`
	}

	RuleInfo struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Event       string `json:"event"`
		DateCreated string `json:"dateCreated"`
		DateUpdated string `json:"dateUpdated"`
		IsActive    bool   `json:"isActive"`
		Ran         string `json:"ran"`
	}
)

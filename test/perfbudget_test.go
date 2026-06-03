package noise

import "os"

// relaxPerfBudgets reports whether wall-clock latency assertions should be
// skipped. Timing is unreliable on shared CI runners (and under the race
// detector), so CI relaxes the budget checks while the functional assertions
// in the same tests still run. Local runs keep the strict budgets.
func relaxPerfBudgets() bool { return os.Getenv("CI") != "" }

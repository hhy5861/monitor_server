package monitor

import "testing"

var (
	formula1 = "(endtime - statustime) > 5"
	formula2 = "((endtime - statustime) > 5) && ((now + 120) < field)"
	formula3 = "((endtime - statustime) > 5) && ((now - 120) < field) && ((now - 120) < field)"
	formula4 = "(now - 120) <= field"
	formula5 = "field > 5"
)

func TestFormulaAnalyze(t *testing.T) {
	FormulaAnalyze(formula5)
}

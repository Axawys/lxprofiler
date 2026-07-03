package detect

import (
	"os/exec"

	"github.com/Axawys/lxprofiler/internal/data"
)

var Score = map[string]int{}
var Reasons = map[string]string{}

func init() {
	for key := range data.Labels {
		Score[key] = 0
		Reasons[key] = ""
	}
}

func Add(key string, pts int, reason string) {
	Score[key] += pts
	if Reasons[key] != "" {
		Reasons[key] += ", "
	}
	Reasons[key] += reason
}

func Has(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

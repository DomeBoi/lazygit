package helpers

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

func (self *RefsHelper) PromptForBranchName(opts types.PromptOpts) {
	onConfirm := opts.HandleConfirm
	opts.HandleConfirm = func(response string) error {
		return onConfirm(SanitizedBranchName(response))
	}
	self.c.Prompt(opts)
}

// SanitizedBranchName will remove all spaces in favor of a dash "-" to meet
// git's branch naming requirement.
func SanitizedBranchName(input string) string {
	return strings.ReplaceAll(input, " ", "-")
}

package helpers

import (
	"errors"
	"regexp"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

func (self *RefsHelper) PromptForBranchName(opts types.PromptOpts) {
	self.promptForBranchName(opts, "")
}

func (self *RefsHelper) PromptForBranchRename(opts types.PromptOpts) {
	self.promptForBranchName(opts, opts.InitialContent)
}

func (self *RefsHelper) promptForBranchName(opts types.PromptOpts, renamedBranch string) {
	onConfirm := opts.HandleConfirm
	opts.PreserveWhitespace = true
	opts.GetFeedback = func(input string) string {
		return self.branchNameFeedback(input, renamedBranch)
	}
	opts.ValidateInput = func(input string) error {
		return self.validateBranchName(input, renamedBranch)
	}
	opts.HandleConfirm = func(response string) error {
		return onConfirm(SanitizedBranchName(response))
	}
	self.c.Prompt(opts)
}

// Match GitHub Desktop's ref-name corrections, including characters that cannot
// be used in refs on Windows. Git still validates the result before submission.
var invalidBranchNameCharacters = regexp.MustCompile(`[\x00-\x20\x7F~^:?*\[\\|"<>]+|@\{|\.\.+|^\.|\.$|\.lock$|/$`)

func SanitizedBranchName(input string) string {
	return strings.TrimLeft(invalidBranchNameCharacters.ReplaceAllString(input, "-"), "-+")
}

func (self *RefsHelper) branchNameError(name string, renamedBranch string) error {
	if name == "" {
		return errors.New(self.c.Tr.BranchNameCannotBeEmpty)
	}
	if name != renamedBranch && lo.ContainsBy(self.c.Model().Branches, func(branch *models.Branch) bool { return branch.Name == name }) {
		return errors.New(utils.ResolvePlaceholderString(self.c.Tr.BranchNameAlreadyExists, map[string]string{"branchName": name}))
	}
	return nil
}

func (self *RefsHelper) validateBranchName(input string, renamedBranch string) error {
	name := SanitizedBranchName(input)
	if err := self.branchNameError(name, renamedBranch); err != nil {
		return err
	}
	if err := self.c.Git().Branch.CheckName(name); err != nil {
		return errors.New(utils.ResolvePlaceholderString(self.c.Tr.InvalidBranchName, map[string]string{"branchName": name}))
	}
	return nil
}

func (self *RefsHelper) branchNameFeedback(input string, renamedBranch string) string {
	if input == "" {
		return ""
	}
	name := SanitizedBranchName(input)
	if err := self.branchNameError(name, renamedBranch); err != nil {
		return style.FgRed.Sprint(err.Error())
	}
	if name == input {
		return ""
	}
	message := self.c.Tr.BranchNameWillBeCreatedAs
	if renamedBranch != "" {
		message = self.c.Tr.BranchNameWillBeRenamedTo
	}
	return style.FgYellow.Sprint(utils.ResolvePlaceholderString(message, map[string]string{"branchName": name}))
}

package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var RenameBranchNameCorrection = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Preview the corrected rename and allow retaining the branch's own name",
	SetupConfig: func(*config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().Focus().Press(keys.Branches.RenameBranch)
		t.ExpectPopup().Prompt().Title(Contains("Enter new branch name")).InitialText(Equals("master"))
		t.Views().PromptFeedback().IsInvisible()
		t.ExpectPopup().Prompt().Title(Contains("Enter new branch name")).Confirm()
		t.Git().CurrentBranchName("master")
		t.Views().Branches().Press(keys.Branches.RenameBranch)
		t.ExpectPopup().Prompt().Title(Contains("Enter new branch name")).Clear().Type("my branch:name")
		t.Views().PromptFeedback().IsVisible().Content(Equals("Will be renamed to 'my-branch-name'."))
		t.ExpectPopup().Prompt().Title(Contains("Enter new branch name")).Confirm()
		t.Git().CurrentBranchName("my-branch-name")
		t.Views().PromptFeedback().IsInvisible()
	},
})

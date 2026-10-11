package worktree

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NewWorktreeNameCorrection = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Correct new branch names in the worktree picker without turning collisions into branch selections",
	SetupConfig: func(*config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.NewBranch("taken-name")
		shell.Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Worktrees().Focus().Press(keys.Universal.New)
		t.ExpectPopup().Prompt().Title(Equals("New worktree for branch")).Type("taken name")
		t.Views().PromptFeedback().IsVisible().Content(Equals("A branch named 'taken-name' already exists."))
		t.ExpectPopup().Prompt().Title(Equals("New worktree for branch")).Confirm()
		t.ExpectToast(Equals("A branch named 'taken-name' already exists."))
		t.Views().Prompt().IsFocused().Content(Equals("taken name"))

		t.ExpectPopup().Prompt().Title(Equals("New worktree for branch")).Clear().Type("new branch:name")
		t.Views().PromptFeedback().IsVisible().Content(Equals("Will be created as 'new-branch-name'."))
		t.ExpectPopup().Prompt().Title(Equals("New worktree for branch")).Confirm()
		t.Views().PromptFeedback().IsInvisible()
		t.ExpectPopup().Menu().Title(Equals("Worktree location")).Confirm()
		t.Views().Worktrees().IsFocused().ContainsLines(Contains("new-branch-name").IsSelected())
		t.Views().Status().Content(Contains("repo(new-branch-name) → new-branch-name"))
	},
})

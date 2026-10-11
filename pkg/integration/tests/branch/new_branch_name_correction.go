package branch

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var NewBranchNameCorrection = NewIntegrationTest(NewIntegrationTestArgs{
	Description: "Preview corrected branch names, retain rejected input, and create the previewed branch",
	Width:       80,
	Height:      24,
	SetupConfig: func(*config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.EmptyCommit("initial commit")
		shell.NewBranch("taken-name")
		shell.Checkout("master")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Branches().Focus().Press(keys.Universal.New)

		t.ExpectPopup().Prompt().Title(Contains("New branch name"))
		t.Views().PromptFeedback().IsInvisible()

		for _, test := range []struct {
			input string
			error string
		}{
			{"  ?* ", "The corrected branch name cannot be empty."},
			{"taken name", "A branch named 'taken-name' already exists."},
			{"feature//name", "'feature//name' is not a valid branch name."},
			{"feature/.hidden", "'feature/.hidden' is not a valid branch name."},
			{"HEAD", "'HEAD' is not a valid branch name."},
		} {
			t.ExpectPopup().Prompt().Title(Contains("New branch name")).Clear().Type(test.input).Confirm()
			t.ExpectToast(Equals(test.error))
			t.Views().Prompt().IsFocused().Content(Equals(test.input))
			t.Git().CurrentBranchName("master")
		}

		t.ExpectPopup().Prompt().Title(Contains("New branch name")).Clear().Type("valid-name")
		t.Views().PromptFeedback().IsInvisible()

		input := "CORE-55098 createWebFile mutation fails when uploading image/x-icon-files"
		corrected := "CORE-55098-createWebFile-mutation-fails-when-uploading-image/x-icon-files"
		t.ExpectPopup().Prompt().Title(Contains("New branch name")).Clear().Type(input)
		t.Views().Prompt().Content(Equals(input))
		t.Views().PromptFeedback().IsVisible().
			Content(Equals("Will be created as '"+corrected+"'.")).
			ContainsViewLines(Contains("Will be created as"), Contains("icon-files")).
			SharesTopBorderWithBottomOf(t.Views().Prompt())
		t.ExpectPopup().Prompt().Title(Contains("New branch name")).Confirm()
		t.Git().CurrentBranchName(corrected)
		t.Views().PromptFeedback().IsInvisible()

		t.Views().Branches().Press(keys.Universal.New)
		t.ExpectPopup().Prompt().Title(Contains("New branch name")).Type("cancel this")
		t.Views().PromptFeedback().IsVisible()
		t.ExpectPopup().Prompt().Title(Contains("New branch name")).Cancel()
		t.Views().PromptFeedback().IsInvisible()
		t.Git().CurrentBranchName(corrected)
	},
})

# Branch names

When creating or renaming a branch, you can type or paste a name containing
spaces or invalid characters. Lazygit keeps your original text in the input
and shows the corrected name below it. Press Enter to use the displayed name.

For example, `feature/fix image:upload` becomes `feature/fix-image-upload`.
Correction follows GitHub Desktop's rules: invalid characters are replaced
with hyphens, while case and ordinary slash-separated branch prefixes are
preserved.

The same preview is available when moving commits to a new branch and when
creating a branch for a worktree. Selecting an existing branch in the
worktree picker keeps that branch's name.

If the corrected name is empty, already exists, or is still invalid according
to Git, the prompt stays open so you can edit it.

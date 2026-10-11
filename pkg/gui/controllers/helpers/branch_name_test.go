package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizedBranchName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"this-is/fine", "this-is/fine"},
		{"CORE-55098 createWebFile mutation fails when uploading image/x-icon-files", "CORE-55098-createWebFile-mutation-fails-when-uploading-image/x-icon-files"},
		{`.this..is\not fine:yo?|is-it`, "this-is-not-fine-yo-is-it"},
		{"hello/", "hello-"},
		{"++but-can-still-keep-the-rest", "but-can-still-keep-the-rest"},
		{"--but-can-still-keep-the-rest", "but-can-still-keep-the-rest"},
		{"foo.lock.lock", "foo.lock-"},
		{"hello\r\nworld", "hello-world"},
		{".first.dot.is.not.ok", "first.dot.is.not.ok"},
		{"branch--name", "branch--name"},
		{"a  b\t\tc", "a-b-c"},
		{"a\x00\x1f\x7fb", "a-b"},
		{"a~^:?*[b", "a-b"},
		{"a@{b", "a-b"},
		{`a"<>|b`, "a-b"},
		{"name.", "name-"},
		{"feature/修正-é", "feature/修正-é"},
		{"+--+", ""},
		{"  ?* ", ""},
		{"", ""},
		{" name ", "name-"},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			assert.Equal(t, test.want, SanitizedBranchName(test.input))
		})
	}
}

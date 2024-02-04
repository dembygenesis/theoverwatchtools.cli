package main

/* ==================================================
	Available CLI commands are enumerated below:
   ==================================================
*/

type command string

const (
	clipGptPreface   command = "clip-gpt-preface"
	clipFileContents command = "clip-file-contents"
	copyFolderAToB   command = "copy-folder-a-to-b"
)

func (c command) string() string {
	return string(c)
}

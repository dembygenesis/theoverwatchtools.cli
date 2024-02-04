package models

type CopyOptions struct {
	// Source is the destination to copy against
	Source string

	// SourceExclusions are the files/folders to exclude copying.
	SourceExclusions []string

	// Destination is the folder to paste the destination files to,
	// default value is it overwrites.
	Destination string

	// WipeDestination removes all the contents of the destination folder
	// before deleting the data.
	WipeDestination bool

	// WipeDestinationExclusions will exclude files from being deleted if
	// WipeDestination is set as true.
	WipeDestinationExclusions []string
}

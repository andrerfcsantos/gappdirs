package gappdirs

import (
	"io"
	"io/fs"
)

// DefaultCreateFilePerm is the default file permission used for file creation operations.
const DefaultCreateFilePerm fs.FileMode = 0o644

// CreateFileOptions contains the options passed on file creation operations.
type CreateFileOptions struct {
	filePerm          fs.FileMode
	overwriteExisting bool
	reader            io.Reader
}

// CreateFileOption allows setting one option on file creation operations.
type CreateFileOption func(*CreateFileOptions)

func defaultCreateFileOptions() CreateFileOptions {
	return CreateFileOptions{
		filePerm: DefaultCreateFilePerm,
	}
}

// WithFilePerm sets the file permission used when creating new files.
//
// If this option is not set a default permission of 0o644 is used.
// Invalid permissions are ignored and the default is used instead.
func WithFilePerm(perm fs.FileMode) CreateFileOption {
	return func(cfg *CreateFileOptions) {
		if cfg == nil {
			return
		}
		if !isValidCreateFilePerm(perm) {
			return
		}
		cfg.filePerm = perm
	}
}

// WithContentsFromReader sets the contents of the new file being created to the contents of the reader.
//
// If this option is not set a new file will be created with no content.
func WithContentsFromReader(reader io.Reader) CreateFileOption {
	return func(cfg *CreateFileOptions) {
		if cfg == nil {
			return
		}
		cfg.reader = reader
	}
}

// WithOverwriteExisting allows overwriting existing files in file creation operations.
//
// If the option WithContentsFromReader is also used, the contents of the reader will be written to the file, overwriting the existing contents.
func WithOverwriteExisting() CreateFileOption {
	return func(cfg *CreateFileOptions) {
		if cfg == nil {
			return
		}
		cfg.overwriteExisting = true
	}
}

func resolveCreateFileOptions(opts []CreateFileOption) CreateFileOptions {
	cfg := defaultCreateFileOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&cfg)
	}
	if !isValidCreateFilePerm(cfg.filePerm) {
		cfg.filePerm = DefaultCreateFilePerm
	}
	return cfg
}

func isValidCreateFilePerm(perm fs.FileMode) bool {
	return perm != 0 && perm&fs.ModeType == 0
}

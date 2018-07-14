// Package repository contains helper methods for working with a Git repo.
package repository

import (
	"bytes"
	"github.com/MichaelMure/git-bug/util"
	"strings"
)

// Repo represents a source code repository.
type Repo interface {
	// GetPath returns the path to the repo.
	GetPath() string

	// GetUserName returns the name the the user has used to configure git
	GetUserName() (string, error)

	// GetUserEmail returns the email address that the user has used to configure git.
	GetUserEmail() (string, error)

	// GetCoreEditor returns the name of the editor that the user has used to configure git.
	GetCoreEditor() (string, error)

	// PullRefs pull git refs from a remote
	PullRefs(remote string, refPattern string, remoteRefPattern string) error

	// PushRefs push git refs to a remote
	PushRefs(remote string, refPattern string) error

	// StoreData will store arbitrary data and return the corresponding hash
	StoreData(data []byte) (util.Hash, error)

	// ReadData will attempt to read arbitrary data from the given hash
	ReadData(hash util.Hash) ([]byte, error)

	// StoreTree will store a mapping key-->Hash as a Git tree
	StoreTree(mapping []TreeEntry) (util.Hash, error)

	// StoreCommit will store a Git commit with the given Git tree
	StoreCommit(treeHash util.Hash) (util.Hash, error)

	// StoreCommit will store a Git commit with the given Git tree
	StoreCommitWithParent(treeHash util.Hash, parent util.Hash) (util.Hash, error)

	// UpdateRef will create or update a Git reference
	UpdateRef(ref string, hash util.Hash) error

	// ListRefs will return a list of Git ref matching the given refspec
	ListRefs(refspec string) ([]string, error)

	// ListCommits will return the list of tree hashes of a ref, in chronological order
	ListCommits(ref string) ([]util.Hash, error)

	// ListEntries will return the list of entries in a Git tree
	ListEntries(hash util.Hash) ([]TreeEntry, error)
}

func prepareTreeEntries(entries []TreeEntry) bytes.Buffer {
	var buffer bytes.Buffer

	for _, entry := range entries {
		buffer.WriteString(entry.Format())
	}

	return buffer
}

func readTreeEntries(s string) ([]TreeEntry, error) {
	splitted := strings.Split(s, "\n")

	casted := make([]TreeEntry, len(splitted))
	for i, line := range splitted {
		if line == "" {
			continue
		}

		entry, err := ParseTreeEntry(line)

		if err != nil {
			return nil, err
		}

		casted[i] = entry
	}

	return casted, nil
}

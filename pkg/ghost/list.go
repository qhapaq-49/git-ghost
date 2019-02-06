package ghost

import (
	"bytes"
	"fmt"
	"git-ghost/pkg/ghost/types"
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
)

// ListOptions represents arg for List func
type ListOptions struct {
	types.WorkingEnvSpec
	*types.ListCommitsBranchSpec
	*types.ListDiffBranchSpec
}

// ListResult contains results of List func

type ListResult struct {
	*types.LocalBaseBranches
	*types.LocalModBranches
}

// List returns ghost branches list per ghost branch type
func List(options ListOptions) (*ListResult, error) {
	log.WithFields(util.ToFields(options)).Debug("list command with")

	res := ListResult{}

	if options.ListCommitsBranchSpec != nil {
		resolved := options.ListCommitsBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo)
		if err != nil {
			return nil, err
		}
		branches.Sort()
		res.LocalBaseBranches = &branches
	}

	if options.ListDiffBranchSpec != nil {
		resolved := options.ListDiffBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo)
		if err != nil {
			return nil, err
		}
		branches.Sort()
		res.LocalModBranches = &branches
	}

	return &res, nil
}

// PrettyString pretty prints ListResult
func (res *ListResult) PrettyString() string {
	// TODO: Make it prettier
	var buffer bytes.Buffer
	if res.LocalBaseBranches != nil {
		buffer.WriteString("Local Base Branches:\n")
		buffer.WriteString("\n")
		for _, branch := range *res.LocalBaseBranches {
			buffer.WriteString(fmt.Sprintf("%s => %s\n", branch.RemoteBaseCommit, branch.LocalBaseCommit))
		}
		buffer.WriteString("\n")
	}
	if res.LocalModBranches != nil {
		buffer.WriteString("Local Mod Branches:\n")
		buffer.WriteString("\n")
		for _, branch := range *res.LocalModBranches {
			buffer.WriteString(fmt.Sprintf("%s -> %s\n", branch.LocalBaseCommit, branch.LocalModHash))
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

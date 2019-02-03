package ghost

import (
	"bytes"
	"fmt"
	"git-ghost/pkg/ghost/git"
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
)

type DeleteOptions struct {
	WorkingEnvSpec
	Prefix string
	*ListCommitsBranchSpec
	*ListDiffBranchSpec
	Dryrun bool
}

type DeleteResult struct {
	LocalBaseBranches LocalBaseBranches
	LocalModBranches  LocalModBranches
}

func Delete(options DeleteOptions) (*DeleteResult, error) {
	log.WithFields(util.ToFields(options)).Debug("delete command with")

	res := DeleteResult{}

	if options.ListCommitsBranchSpec != nil {
		resolved := options.ListCommitsBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo, options.Prefix)
		if err != nil {
			return nil, err
		}
		res.LocalBaseBranches = branches
	}

	if options.ListDiffBranchSpec != nil {
		resolved := options.ListDiffBranchSpec.Resolve(options.SrcDir)
		branches, err := resolved.GetBranches(options.GhostRepo, options.Prefix)
		if err != nil {
			return nil, err
		}
		res.LocalModBranches = branches
	}

	workingEnv, err := options.WorkingEnvSpec.initialize()
	if err != nil {
		return nil, err
	}
	defer workingEnv.clean()

	deleteBranches := func(branches []GhostBranch, dryrun bool) error {
		var branchNames []string
		for _, branch := range branches {
			branchNames = append(branchNames, branch.BranchName())
		}
		log.WithFields(log.Fields{
			"branches": branchNames,
		}).Info("Delete branch")
		if dryrun {
			return nil
		}
		return git.DeleteRemoteBranches(workingEnv.GhostDir, branchNames...)
	}

	if len(res.LocalBaseBranches) > 0 {
		res.LocalBaseBranches.Sort()
		err := deleteBranches(res.LocalBaseBranches.AsGhostBranches(), options.Dryrun)
		if err != nil {
			return nil, err
		}
	}

	if len(res.LocalModBranches) > 0 {
		res.LocalModBranches.Sort()
		err := deleteBranches(res.LocalModBranches.AsGhostBranches(), options.Dryrun)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (res *DeleteResult) PrettyString() string {
	// TODO: Make it prettier
	var buffer bytes.Buffer
	if len(res.LocalBaseBranches) > 0 {
		buffer.WriteString("Deleted Local Base Branches:\n")
		buffer.WriteString("\n")
	}
	for _, branch := range res.LocalBaseBranches {
		buffer.WriteString(fmt.Sprintf("%s => %s\n", branch.RemoteBaseCommit, branch.LocalBaseCommit))
	}
	if len(res.LocalBaseBranches) > 0 {
		buffer.WriteString("\n")
	}
	if len(res.LocalModBranches) > 0 {
		buffer.WriteString("Deleted Local Mod Branches:\n")
		buffer.WriteString("\n")
	}
	for _, branch := range res.LocalModBranches {
		buffer.WriteString(fmt.Sprintf("%s -> %s\n", branch.LocalBaseCommit, branch.LocalModHash))
	}
	if len(res.LocalModBranches) > 0 {
		buffer.WriteString("\n")
	}
	return buffer.String()
}

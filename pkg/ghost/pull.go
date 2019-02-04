package ghost

import (
	"git-ghost/pkg/util"

	log "github.com/Sirupsen/logrus"
)

type PullOptions struct {
	WorkingEnvSpec
	LocalBaseBranchSpec
	PullableLocalModBranchSpec
}

type PullCommitsOptions struct {
	WorkingEnvSpec
	LocalBaseBranchSpec
}
type PullDiffOptions struct {
	WorkingEnvSpec
	PullableLocalModBranchSpec
}

func PullCommits(options PullCommitsOptions, workingEnv *WorkingEnv) error {
	log.WithFields(util.ToFields(options)).Debug("pull commits command with")

	we, initialized, err := initializeWorkingEnvIfRequired(options.WorkingEnvSpec, workingEnv)
	if err != nil {
		return err
	}
	if initialized {
		defer we.clean()
	}

	pulledBranch, err := options.LocalBaseBranchSpec.PullBranch(*we)
	if err != nil {
		return err
	}

	return pulledBranch.Apply(*we)
}

func PullDiff(options PullDiffOptions, workingEnv *WorkingEnv) error {
	log.WithFields(util.ToFields(options)).Debug("pull diff command with")
	we, initialized, err := initializeWorkingEnvIfRequired(options.WorkingEnvSpec, workingEnv)
	if err != nil {
		return err
	}
	if initialized {
		defer we.clean()
	}

	pulledBranch, err := options.PullableLocalModBranchSpec.PullBranch(*we)
	if err != nil {
		return err
	}

	return pulledBranch.Apply(*we)
}

func Pull(options PullOptions) error {
	log.WithFields(util.ToFields(options)).Debug("pull command with")
	we, err := options.WorkingEnvSpec.initialize()
	if err != nil {
		return err
	}
	defer we.clean()

	err = PullCommits(PullCommitsOptions{
		WorkingEnvSpec:      options.WorkingEnvSpec,
		LocalBaseBranchSpec: options.LocalBaseBranchSpec,
	}, we)
	if err != nil {
		return err
	}

	err = PullDiff(PullDiffOptions{
		WorkingEnvSpec:             options.WorkingEnvSpec,
		PullableLocalModBranchSpec: options.PullableLocalModBranchSpec,
	}, we)
	if err != nil {
		return err
	}

	return nil
}

func initializeWorkingEnvIfRequired(spec WorkingEnvSpec, we *WorkingEnv) (*WorkingEnv, bool, error) {
	if we == nil {
		we, err := spec.initialize()
		if err != nil {
			return nil, false, err
		}
		return we, true, nil
	}
	return we, false, nil
}

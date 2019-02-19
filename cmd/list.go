package cmd

import (
	"fmt"
	"git-ghost/pkg/ghost"
	"git-ghost/pkg/ghost/types"
	"git-ghost/pkg/util/errors"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(NewListCommand())
}

var outputTypes = []string{"only-from", "only-to"}
var regexpOutputPattern = regexp.MustCompile("^(|" + strings.Join(outputTypes, "|") + ")$")

type listFlags struct {
	hashFrom  string
	hashTo    string
	noHeaders bool
	output    string
}

func NewListCommand() *cobra.Command {
	var (
		listFlags listFlags
	)

	var command = &cobra.Command{
		Use:   "list",
		Short: "list ghost branches of diffs.",
		Long:  "list ghost branches of diffs.",
		Args:  cobra.NoArgs,
		Run:   runListDiffCommand(&listFlags),
	}
	command.AddCommand(&cobra.Command{
		Use:   "commits",
		Short: "list ghost branches of commits.",
		Long:  "list ghost branches of commits.",
		Args:  cobra.NoArgs,
		Run:   runListCommitsCommand(&listFlags),
	})
	command.AddCommand(&cobra.Command{
		Use:   "diff",
		Short: "list ghost branches of diffs.",
		Long:  "list ghost branches of diffs.",
		Args:  cobra.NoArgs,
		Run:   runListDiffCommand(&listFlags),
	})
	command.AddCommand(&cobra.Command{
		Use:   "all",
		Short: "list ghost branches of all types.",
		Long:  "list ghost branches of all types.",
		Args:  cobra.NoArgs,
		Run:   runListAllCommand(&listFlags),
	})
	command.PersistentFlags().StringVar(&listFlags.hashFrom, "from", "", "commit or diff hash to which ghost branches are listed.")
	command.PersistentFlags().StringVar(&listFlags.hashTo, "to", "", "commit or diff hash from which ghost branches are listed.")
	command.PersistentFlags().BoolVar(&listFlags.noHeaders, "no-headers", false, "When using the default, only-from or only-to output format, don't print headers (default print headers).")
	command.PersistentFlags().StringVarP(&listFlags.output, "output", "o", "", "Output format. One of: only-from|only-to")
	return command
}

func runListCommitsCommand(flags *listFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		opts := ghost.ListOptions{
			WorkingEnvSpec: types.WorkingEnvSpec{
				SrcDir:          globalOpts.srcDir,
				GhostWorkingDir: globalOpts.ghostWorkDir,
				GhostRepo:       globalOpts.ghostRepo,
			},
			ListCommitsBranchSpec: &types.ListCommitsBranchSpec{
				Prefix:   globalOpts.ghostPrefix,
				HashFrom: flags.hashFrom,
				HashTo:   flags.hashTo,
			},
		}

		res, err := ghost.List(opts)
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		fmt.Printf(res.PrettyString(!flags.noHeaders, flags.output))
	}
}

func runListDiffCommand(flags *listFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		opts := ghost.ListOptions{
			WorkingEnvSpec: types.WorkingEnvSpec{
				SrcDir:          globalOpts.srcDir,
				GhostWorkingDir: globalOpts.ghostWorkDir,
				GhostRepo:       globalOpts.ghostRepo,
			},
			ListDiffBranchSpec: &types.ListDiffBranchSpec{
				Prefix:   globalOpts.ghostPrefix,
				HashFrom: flags.hashFrom,
				HashTo:   flags.hashTo,
			},
		}

		res, err := ghost.List(opts)
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		fmt.Printf(res.PrettyString(!flags.noHeaders, flags.output))
	}
}

func runListAllCommand(flags *listFlags) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := flags.validate()
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		opts := ghost.ListOptions{
			WorkingEnvSpec: types.WorkingEnvSpec{
				SrcDir:          globalOpts.srcDir,
				GhostWorkingDir: globalOpts.ghostWorkDir,
				GhostRepo:       globalOpts.ghostRepo,
			},
			ListCommitsBranchSpec: &types.ListCommitsBranchSpec{
				Prefix:   globalOpts.ghostPrefix,
				HashFrom: flags.hashFrom,
				HashTo:   flags.hashTo,
			},
			ListDiffBranchSpec: &types.ListDiffBranchSpec{
				Prefix:   globalOpts.ghostPrefix,
				HashFrom: flags.hashFrom,
				HashTo:   flags.hashTo,
			},
		}

		res, err := ghost.List(opts)
		if err != nil {
			errors.LogErrorWithStack(err)
			os.Exit(1)
		}
		fmt.Printf(res.PrettyString(!flags.noHeaders, flags.output))
	}
}

func (flags listFlags) validate() errors.GitGhostError {
	if !regexpOutputPattern.MatchString(flags.output) {
		return errors.Errorf("output must be one of %v", outputTypes)
	}
	return nil
}

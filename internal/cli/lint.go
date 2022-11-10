package cli

import (
	"errors"
	"fmt"
	"github.com/catppuccin/cli/schema"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(lintCmd)
}

var lintCmd = &cobra.Command{
	Use:   "lint .catppuccin.yaml",
	Short: "Lint a port config",
	Long:  "Lints a port configuration against the schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		if err := linter(args[0]); err != nil {
			var errs schema.ResultErrors
			if errors.As(err, &errs) {
				fmt.Println("Port config is invalid!")
				fmt.Println()
				for _, e := range errs {
					fmt.Println(e)
				}
				return nil
			}
			return err
		}
		fmt.Println("Port config is valid!")
		return nil
	},
}

func linter(filePath string) error {
	fi, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fi.Close()

	return schema.Lint(fi)
}

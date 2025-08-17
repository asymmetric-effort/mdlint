// Copyright 2024 MdLint Authors

package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/asymmetric-effort/mdlint/internal/config"
	"github.com/asymmetric-effort/mdlint/internal/engine"
	"github.com/asymmetric-effort/mdlint/internal/formatter"
	"github.com/spf13/cobra"
)

const version = "0.1.0"

func main() {
	var (
		cfgPath     string
		quiet       bool
		formatFlag  string
		listRules   bool
		showVersion bool
	)
	exitCode := 0
	rootCmd := &cobra.Command{
		Use:          "mdlint [files...]",
		Short:        "mdlint lints Markdown files",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if quiet {
				log.SetOutput(io.Discard)
			}
			if showVersion {
				fmt.Fprintf(cmd.OutOrStdout(), "mdlint %s\n", version)
				return nil
			}
			if listRules {
				fmt.Fprintln(cmd.OutOrStdout(), "MD9000 TODO found")
				return nil
			}
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}
			if formatFlag != "" {
				cfg.Format = formatFlag
			}
			eng := engine.Engine{}
			fs, err := eng.Run(args)
			if err != nil {
				return err
			}
			if len(fs) == 0 {
				return nil
			}
			out, err := formatter.Format(fs, cfg.Format)
			if err != nil {
				return err
			}
			if out != "" {
				fmt.Fprint(cmd.OutOrStdout(), out)
			}
			if len(fs) > 0 {
				exitCode = 1
			}
			return nil
		},
	}
	rootCmd.Flags().StringVar(&cfgPath, "config", "", "config file")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "suppress logs")
	rootCmd.Flags().StringVar(&formatFlag, "format", "", "output format")
	rootCmd.Flags().BoolVar(&listRules, "list-rules", false, "list supported rules")
	rootCmd.Flags().BoolVar(&showVersion, "version", false, "print version")
	rootCmd.SilenceErrors = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	os.Exit(exitCode)
}

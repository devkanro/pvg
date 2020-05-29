package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var rootCmd = &cobra.Command{
	Use:     "pvg",
	Short:   "pvg is tool for convert pixel art to SVG.",
	Long:    `Convert pixel art bitmap to SVG without tracing, just keep all pixel exist.`,
	Version: "v1.0",
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		output := input
		if len(args) == 2 {
			output = args[1]
		}

		transparent, err := parseHexColor(transparentColor)
		if err != nil {
			return fmt.Errorf("Wrong transparent color '%s', it must be a valid hex color value ", transparentColor)
		}

		stat, err := os.Stat(input)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			err = handleDir(input, output, transparent)
		} else {
			basename := filepath.Base(input)
			basename = basename[:len(basename)-len(filepath.Ext(basename))]
			if input == output {
				output = filepath.Join(filepath.Dir(input), basename+".svg")
			}
			stat, _ := os.Stat(output)
			if stat == nil {
				if filepath.Ext(output) == ".svg" {
					err = os.MkdirAll(filepath.Dir(output), os.ModeDir)
					if err != nil {
						return err
					}
				} else {
					err = os.MkdirAll(output, os.ModeDir)
					if err != nil {
						return err
					}
					output = filepath.Join(output, basename+".svg")
				}
			} else if stat.IsDir() {
				output = filepath.Join(output, basename+".svg")
			}
			err = handleFile(input, output, transparent)
		}
		return err
	},
}

var transparentColor string
var parallel bool

func init() {
	rootCmd.Flags().StringVarP(&transparentColor, "transparent", "t", "#FFFF00FF",
		"transparent color for converting, default as purple.")
	rootCmd.Flags().BoolVarP( &parallel, "parallel", "p", true,
		"parallel process in folder mode.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

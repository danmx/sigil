// +build ignore

package cmd

import (
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	outDir = path.Join("docs", "man")

	// genDocCmd represents the gendoc command
	genDocCmd = &cobra.Command{
		Use:   "gendoc",
		Short: "Generate the documentation in Markdown",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				return err
			}
			docDir := path.Join(dir, outDir)
			log.WithFields(log.Fields{
				"docDir":     docDir,
				"output-dir": cfg.GetString("output-dir"),
			}).Debug("Path to the documentation")
			if err = doc.GenMarkdownTree(rootCmd, docDir); err != nil {
				return err
			}
			return nil
		},
		TraverseChildren:  false,
		DisableAutoGenTag: true,
	}
)

func init() {
	rootCmd.AddCommand(genDocCmd)

	genDocCmd.Flags().StringVar(&outDir, "output-dir", outDir, "specify output directory")
}

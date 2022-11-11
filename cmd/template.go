package cmd

import (
	"fmt"
	"os"

	"github.com/jbchouinard/wmt/editor"
	"github.com/jbchouinard/wmt/env"
	"github.com/jbchouinard/wmt/errored"
	"github.com/jbchouinard/wmt/template"
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Template commands",
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Args:  cobra.ExactArgs(0),
	Short: "List templates",
	Run: func(cmd *cobra.Command, args []string) {
		templates := template.ListTemplates()
		for _, tmpl := range templates {
			fmt.Println(tmpl)
		}
	},
}

var templateNewCmd = &cobra.Command{
	Use:   "add name",
	Args:  cobra.ExactArgs(1),
	Short: "Create a new template",
	Long: `Create a new template. Variables like {{.foo}} are substituted.
	
See https://pkg.go.dev/text/template for details on template syntax.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		tmpl, err := template.CreateTemplate(name, newTemplateHtml)
		errored.Check(err, "create template")
		filename := tmpl.Path()
		err = editor.Edit(filename, true)
		errored.Check(err, "edit template")
		err = tmpl.Parse()
		errored.Check(err, "validate")
	},
}

var templateEditCmd = &cobra.Command{
	Use:   "edit name",
	Args:  cobra.ExactArgs(1),
	Short: "Edit an existing template",
	Long: `Edit an existing template. Variables like {{.foo}} are substituted.
	
See https://pkg.go.dev/text/template for details on template syntax.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		tmpl, err := template.SelectTemplate(name)
		errored.Check(err, "load template")
		err = editor.Edit(tmpl.Path(), true)
		errored.Check(err, "edit template")
		err = tmpl.Parse()
		errored.Check(err, "validate")
	},
}

var templateDeleteCmd = &cobra.Command{
	Use:   "delete name",
	Args:  cobra.ExactArgs(1),
	Short: "Delete template",
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		tmpl, err := template.SelectTemplate(name)
		errored.Check(err, "load template")
		filename := tmpl.Path()
		errored.Check(os.Remove(filename), "delete template file")
		tmpl.Delete()
	},
}

var templateEvalCmd = &cobra.Command{
	Use:   "eval name",
	Args:  cobra.ExactArgs(1),
	Short: "Evaluate template",
	Run: func(cmd *cobra.Command, args []string) {
		params, err := env.ParseKVs(evalTemplateParams)
		errored.Check(err, "parse params")

		name := args[0]
		tmpl, err := template.SelectTemplate(name)
		errored.Check(err, "load template")
		bytes, err := tmpl.Eval(env.Current, params)
		errored.Check(err, "template eval")
		fmt.Printf("%s", bytes)
	},
}

var newTemplateHtml bool
var evalTemplateParams []string

func init() {
	templateNewCmd.Flags().BoolVar(&newTemplateHtml, "html", false, "is HTML template")
	templateCmd.AddCommand(templateNewCmd)

	templateCmd.AddCommand(templateListCmd)

	templateCmd.AddCommand(templateEditCmd)

	templateCmd.AddCommand(templateDeleteCmd)

	templateEvalCmd.Flags().StringSliceVarP(&evalTemplateParams, "param", "p", make([]string, 0), "parameter like foo=bar")
	templateCmd.AddCommand(templateEvalCmd)

	rootCmd.AddCommand(templateCmd)
}

package cmd

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/jbchouinard/wmt/editor"
	"github.com/jbchouinard/wmt/env"
	"github.com/jbchouinard/wmt/errored"
	"github.com/jbchouinard/wmt/http"
	"github.com/jbchouinard/wmt/template"
	"github.com/spf13/cobra"
)

var reqCmd = &cobra.Command{
	Use:   "req",
	Short: "Edit and do HTTP requests",
}

var reqAddCmd = &cobra.Command{
	Use:   "add name method path",
	Short: "Add or update saved HTTP request definition",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		req, err := http.LoadRequestDefinition(args[0])
		if err == nil {
			req.Method = args[1]
			req.URL = args[2]
		} else if err == sql.ErrNoRows {
			req = http.NewRequestDefinition(args[0], args[1], args[2])
		} else {
			errored.Fatal(err.Error())
		}
		if editTemplate {
			if req.Template == nil {
				req.Template, err = template.CreateTemplate(fmt.Sprintf("request.%s.body", req.Name), false)
				errored.Check(err, "")
			}
			errored.Check(editor.Edit(req.Template.Path(), true), "")
			fmt.Println(req.Template.Id())
		}
		errored.Check(req.Save(), "")
	},
}

var reqDeleteCmd = &cobra.Command{
	Use:   "delete name",
	Short: "Delete HTTP request definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		existing, err := http.LoadRequestDefinition(args[0])
		if err == nil {
			errored.Check(existing.Delete(), "")
		} else if err == sql.ErrNoRows {
			errored.Fatal("no request with this name")
		} else {
			errored.Fatal(err.Error())
		}
	},
}

var reqShowCmd = &cobra.Command{
	Use:   "show name",
	Short: "Show saved HTTP request definition",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req, err := http.LoadRequestDefinition(args[0])
		if err == nil {
			fmt.Println(req.Details())
		} else if err == sql.ErrNoRows {
			errored.Fatal("no request with this name")
		} else {
			errored.Fatal(err.Error())
		}
	},
}

var reqListCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved HTTP request definitions",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		names := http.ListRequestDefinitions()
		for _, name := range names {
			req, err := http.LoadRequestDefinition(name)
			errored.Check(err, "")
			fmt.Println(req)
		}
	},
}

var reqHeaderCmd = &cobra.Command{
	Use:   "header name key value",
	Short: "Set header for request definition",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[1]
		value := args[2]
		req, err := http.LoadRequestDefinition(args[0])
		if err == nil {
			if value == "" || value == "_" {
				delete(req.Headers, key)
			} else {
				req.Headers[key] = value
			}
			errored.Check(req.Save(), "")
		} else if err == sql.ErrNoRows {
			errored.Fatal("no request with this name")
		} else {
			errored.Fatal(err.Error())
		}
	},
}

var reqQueryCmd = &cobra.Command{
	Use:   "query name key value",
	Short: "Set query param for request definition",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[1]
		value := args[2]
		req, err := http.LoadRequestDefinition(args[0])
		if err == nil {
			if value == "_" || value == "" {
				delete(req.Query, key)
			} else {
				req.Query[key] = value
			}
			errored.Check(req.Save(), "")
		} else if err == sql.ErrNoRows {
			errored.Fatal("no request with this name")
		} else {
			errored.Fatal(err.Error())
		}
	},
}

var reqWriteCmd = &cobra.Command{
	Use:   "write name",
	Short: "Write HTTP request to stdout",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reqDef, err := http.LoadRequestDefinition(args[0])
		if err == sql.ErrNoRows {
			errored.Fatal("no request with this name")
		} else if err != nil {
			errored.Fatal(err.Error())
		}
		ps, err := env.ParseKVs(params)
		errored.Check(err, "")
		reqDef.Eval(env.Current, ps).Request().Write(os.Stdout)
	},
}

var reqDoCmd = &cobra.Command{
	Use:   "do name",
	Short: "Execute saved HTTP request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reqDef, err := http.LoadRequestDefinition(args[0])
		if err == sql.ErrNoRows {
			errored.Fatal("no request with this name")
		} else if err != nil {
			errored.Fatal(err.Error())
		}
		ps, err := env.ParseKVs(params)
		errored.Check(err, "")
		request := reqDef.Eval(env.Current, ps).Request()
		client := http.MakeClient()
		response, err := client.Do(request)
		errored.Check(err, "")
		fmt.Println(response)
	},
}

var reqGetCmd = &cobra.Command{
	Use:   "get url",
	Short: "Execute HTTP GET request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		doRequest("GET", args[0])
	},
}

var reqPostCmd = &cobra.Command{
	Use:   "post url",
	Short: "Execute HTTP POST request",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		doRequest("POST", args[0])
	},
}

var verbose bool
var params []string
var headers []string
var query []string
var editTemplate bool

func doRequest(method string, url string) {
	reqDef := http.NewRequestDefinition("<anonymous>", method, url)

	paramKvs, err := env.ParseKVs(params)
	errored.Check(err, "")

	headerKvs, err := env.ParseKVs(headers)
	errored.Check(err, "")
	env.AddKVs(reqDef.Headers, headerKvs)

	queryKvs, err := env.ParseKVs(query)
	errored.Check(err, "")
	env.AddKVs(reqDef.Query, queryKvs)

	client := http.MakeClient()
	response, err := client.Do(reqDef.Eval(env.Current, paramKvs).Request())
	errored.Check(err, "")
	if verbose {
		response.Write(os.Stdout)
	} else {
		fmt.Println(response.Status)
		body, err := io.ReadAll(response.Body)
		errored.Check(err, "")
		fmt.Println(string(body))
	}
}

func doRequestFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	cmd.Flags().StringSliceVarP(&params, "param", "p", make([]string, 0), "template parameter")
	cmd.Flags().StringSliceVar(&headers, "header", make([]string, 0), "header")
	cmd.Flags().StringSliceVarP(&query, "query", "q", make([]string, 0), "query parameter")
}

func init() {
	reqAddCmd.Flags().BoolVarP(&editTemplate, "template", "t", false, "edit request body template")
	reqCmd.AddCommand(reqAddCmd)

	reqCmd.AddCommand(reqDeleteCmd)
	reqCmd.AddCommand(reqShowCmd)
	reqCmd.AddCommand(reqListCmd)
	reqCmd.AddCommand(reqHeaderCmd)
	reqCmd.AddCommand(reqQueryCmd)

	reqWriteCmd.Flags().StringSliceVarP(&params, "param", "p", make([]string, 0), "template parameter")
	reqCmd.AddCommand(reqWriteCmd)

	reqDoCmd.Flags().StringSliceVarP(&params, "param", "p", make([]string, 0), "template parameter")
	reqCmd.AddCommand(reqDoCmd)

	doRequestFlags(reqGetCmd)
	reqCmd.AddCommand(reqGetCmd)
	doRequestFlags(reqPostCmd)
	reqCmd.AddCommand(reqPostCmd)

	rootCmd.AddCommand(reqCmd)
}

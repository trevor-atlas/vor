package main

import (
	"fmt"
	// "os"

	color "github.com/fatih/color"
	"github.com/trevor-atlas/vor/src/utils"
	cmd "github.com/trevor-atlas/vor/src/cmd"
)

type App struct {
	version      string
	name         string
	description  string
	cmds         [][]string
	requiredEnvs []string
	DEFAULTS     defaults
}

type defaults struct {
	leftPad int
}

func (app *App) handleMissingEnv() {
	errors := []string{}
	for i := range app.requiredEnvs {
		env := app.requiredEnvs[i]
		if utils.GetEnv(env) == "" {
			errors = append(errors, env)
		}
	}
	if len(errors) > 0 {
		fmt.Println(app.name + " error(s):")
		for i := range errors {
			color.Red(utils.LeftPad(errors[i]+" missing", " ", app.DEFAULTS.leftPad))
		}
		fmt.Println()
	}

}

func (app *App) handleNoArgs(cliArgs *[]string) {
	if len(*cliArgs) == 0 {
		cyan := color.New(color.FgCyan).PrintfFunc()
		color.Blue(app.name + " – " + app.description)
		for i := range app.cmds {
			name := utils.LeftPad(app.cmds[i][0], " ", app.DEFAULTS.leftPad)
			desc := app.cmds[i][1]
			cyan(name)
			fmt.Print(": ", utils.LeftPad(desc+"\n", " ", 16-len(name)))
		}
		return
	}
}

func main() {
	// app := App{
	// 	version:     "0.0.1",
	// 	name:        "Vör",
	// 	description: "Jira & Git made simple",
	// 	cmds: [][]string{
	// 		{"branch", "create a branch for a given jira issue number"},
	// 		{"pull-request", "create a PR in github with your current branch"},
	// 		{"issues", "list jira issues assigned to me"},
	// 		{"issue", "get details for a specific issue number"},
	// 		{"review", "open github review for a specific issue number"},
	// 	},
	// 	requiredEnvs: []string{"JIRA_API_KEY", "GITHUB_API_KEY"},
	// 	DEFAULTS: defaults{
	// 		leftPad: 4,
	// 	},
	// }

	// `os.Args` provides access to raw command-line
	// arguments. Note that the first value in this slice
	// is the path to the program, and `os.Args[1:]`
	// holds the arguments to the program.
	// You can get individual args with normal indexing.
	// cliArgs := os.Args[1:]

	// app.handleMissingEnv()
	// app.handleNoArgs(&cliArgs)
	cmd.Execute()
}

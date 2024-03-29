package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/Bananenpro/cli"
	"github.com/spf13/cobra"
)

type doctorRule interface {
	Check() bool
	ErrMessage() string
	SuccessMessage() string
}

type doctorRuleInactive struct{}

func (d doctorRuleInactive) Check() bool {
	return true
}

func (d doctorRuleInactive) ErrMessage() string {
	return ""
}

func (d doctorRuleInactive) SuccessMessage() string {
	return ""
}

type doctorRuleTool struct {
	names   []string
	message string
}

func newDoctorRuleTool(message string, names ...string) doctorRule {
	return doctorRuleTool{
		names:   names,
		message: message,
	}
}

func newDoctorRuleToolWithCondition(message string, active bool, names ...string) doctorRule {
	if !active {
		return doctorRuleInactive{}
	}
	return doctorRuleTool{
		names:   names,
		message: message,
	}
}

func (d doctorRuleTool) Check() bool {
	for _, n := range d.names {
		if _, err := exec.LookPath(n); err == nil {
			return true
		}
	}
	return false
}

func (d doctorRuleTool) ErrMessage() string {
	return d.message
}

func (d doctorRuleTool) SuccessMessage() string {
	for _, n := range d.names {
		if _, err := exec.LookPath(n); err == nil {
			return fmt.Sprintf("`%s` is installed.", n)
		}
	}
	panic("SuccessMessage() called but tool is not installed")
}

var installDir = "/usr/local/bin"

func init() {
	if runtime.GOOS == "windows" {
		installDir = os.Getenv("HOME") + "\\AppData\\Local\\Programs\\codegame-cli"
	}
}

type doctorCategory struct {
	name  string
	rules []doctorRule
}

var doctorRules = []doctorCategory{
	{name: "CLI", rules: []doctorRule{
		newDoctorRuleTool("`codegame` is not in PATH. If you have installed codegame-cli in a custom install directory, make sure to add it to the PATH environment variable. Otherwise, manually add "+installDir+" to the PATH environment variable.", "codegame"),
		newDoctorRuleToolWithCondition("Either curl or wget must be installed to use `codegame upgrade`.", runtime.GOOS != "windows", "curl", "wget"),
	}},
	{name: "C#", rules: []doctorRule{
		newDoctorRuleTool("`dotnet` must be installed to develop CodeGame applications using C#. Install it from https://dotnet.microsoft.com/en-us/download.", "dotnet"),
	}},
	{name: "Go", rules: []doctorRule{
		newDoctorRuleTool("`go` must be installed to develop CodeGame applications using the Go programming language. Install it from https://go.dev.", "go"),
	}},
	{name: "Java", rules: []doctorRule{
		newDoctorRuleTool("`java` must be installed to develop CodeGame applications using Java. Install it from https://adoptium.net.", "java"),
		newDoctorRuleTool("`mvn` must be installed to develop CodeGame applications using Java. Download it from https://maven.apache.org/download.cgi and follow the instructions at https://maven.apache.org/install.html.", "mvn"),
	}},
	{name: "JavaScript", rules: []doctorRule{
		newDoctorRuleTool("`npm` must be installed to develop CodeGame applications using JavaScript or TypeScript. Install it from https://nodejs.org.", "npm"),
		newDoctorRuleTool("`npx` must be installed to run TypeScript or browser based CodeGame applications. Install it using `npm install -g npx`.", "npx"),
	}},
}

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check for missing dependencies and misconfigurations.",
	Run: func(_ *cobra.Command, _ []string) {
		for _, category := range doctorRules {
			cli.PrintColor(cli.Cyan, "%s:", category.name)
			for _, r := range category.rules {
				if _, ok := r.(doctorRuleInactive); ok {
					continue
				}
				if r.Check() {
					cli.PrintColor(cli.Green, "  √ %s", r.SuccessMessage())
				} else {
					cli.PrintColor(cli.Red, "  x %s", r.ErrMessage())
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

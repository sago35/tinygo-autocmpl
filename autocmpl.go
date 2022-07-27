package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

// flagCompleteMap specifies which flags to autocomplete.
var flagCompleteMap = map[string][]string{
	"bench":         {},
	"benchtime":     {},
	"c":             nil,
	"cflags":        {},
	"cpuprofile":    {},
	"deps":          nil,
	"dumpssa":       nil,
	"gc":            {"none", "leaking", "extalloc", "conservative"},
	"heap-size":     {},
	"json":          nil,
	"ldflags":       {},
	"llvm-features": nil,
	"no-debug":      nil,
	"o":             {},
	"ocd-commands":  {},
	"ocd-output":    nil,
	"ocd-verify":    nil,
	"opt":           {"0", "1", "2", "s", "z"},
	"p":             {},
	"panic":         {"print", "trap"},
	"port":          {},
	"print-allocs":  {},
	"print-stacks":  nil,
	"printir":       nil,
	"programmer":    validProgrammers,
	"run":           {},
	"scheduler":     {"none", "tasks", "asyncify", "coroutines"},
	"serial":        {"none", "uart", "usb"},
	"size":          {"none", "short", "full"},
	"short":         nil,
	"tags":          {},
	"target":        validTargets,
	"test":          nil,
	"verifyir":      nil,
	"wasm-abi":      {"generic", "js"},
	"work":          nil,
	"x":             nil,
	"v":             nil,
}

// validTargets is a list of completion targets for -target. It can be overridden by arguments.
var (
	validProgrammers []string
	validTargets     []string
	validCommands    = []string{
		"build",
		"build-library",
		"clang",
		"clean",
		"env",
		"flash",
		"gdb",
		"help",
		"info",
		"ld.lld",
		"list",
		"lldb",
		"run",
		"targets",
		"test",
		"version",
		"wasm-ld",
	}
)

const completionScriptBashStr = `
_tinygo_autocmpl_bash_autocomplete() {
    local cur prev opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    opts=$( tinygo-autocmpl %s -- ${COMP_WORDS[@]:1:$COMP_CWORD} )
    if [ "${opts}" = "" ]; then
        compopt -o filenames
        COMPREPLY=( $(compgen -f -- "${cur}" ) )
    else
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    fi
    return 0
}
complete -F _tinygo_autocmpl_bash_autocomplete tinygo
`

const completionScriptZshStr = `#compdef tinygo

autoload -U compinit && compinit
autoload -U bashcompinit && bashcompinit

_tinygo_autocmpl_bash_autocomplete() {
    local cur prev opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    opts=$( tinygo-autocmpl %s -- ${COMP_WORDS[@]:1:$COMP_CWORD} )
    if [ "${opts}" = "" ]; then
        COMPREPLY=( $(compgen -f -- "${cur}" ) )
    else
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    fi
    return 0
}
complete -o nospace -F _tinygo_autocmpl_bash_autocomplete tinygo
`

func handleCompletionScriptBash(listPath string) {
	targets := ""
	if listPath != "" {
		targets = fmt.Sprintf("--targets=%q", listPath)
	}
	fmt.Printf(completionScriptBashStr, targets)
	return
}

func handleCompletionScriptZsh(listPath string) {
	targets := ""
	if listPath != "" {
		targets = fmt.Sprintf("--targets=%q", listPath)
	}
	fmt.Printf(completionScriptZshStr, targets)
	return
}

const completionScriptClinkStr = `
local parser = clink.arg.new_parser

local tinygo_targets_parser = parser({
    %s,
})

-- function getTargets(dir)
--     local handle = io.popen("dir /b "..dir.."\\*.json", "r")
--     local content = handle:read("*all")
--     handle:close()
-- 
--     local t = {}
--     i = 1
-- 
--     for s in string.gmatch(content, "([^\n]+).json") do
--         t[i] = s
--         i = i + 1
--     end
-- 
--     return t
-- end

-- local tinygo_targets_parser = parser(getTargets("C:\\tinygo\\targets"))


local tinygo_flag_parser = parser(
    %s
    )

local tinygo_parser = parser({
    "build"..tinygo_flag_parser,
    "run"..tinygo_flag_parser,
    "test"..tinygo_flag_parser,
    "flash"..tinygo_flag_parser,
    "gdb"..tinygo_flag_parser,
    "env"..tinygo_flag_parser,
    "list"..tinygo_flag_parser,
    "clean"..tinygo_flag_parser,
    "help"..tinygo_flag_parser,
    })

clink.arg.register_parser("tinygo", tinygo_parser)
`

func init() {
	targets, err := getTargetsFromTinygoTargets()
	if err != nil {
		log.Fatal(err)
	}
	validTargets = targets

	programmers, err := getProgrammers()
	if err != nil {
		log.Fatal(err)
	}
	validProgrammers = programmers

	flagCompleteMap["programmer"] = validProgrammers
	flagCompleteMap["target"] = validTargets
}

func handleCompletionScriptClink(listPath string) {
	targets := []string{}
	for _, t := range validTargets {
		targets = append(targets, fmt.Sprintf("%q", t))
	}

	flags := []string{}
	for _, f := range getFlagCompletion() {
		m, ok := flagCompleteMap[f]
		if !ok {
			panic("panic")
		}

		if 0 < len(m) {
			p := []string{}
			for _, c := range m {
				p = append(p, fmt.Sprintf("%q", c))
			}

			flags = append(flags, fmt.Sprintf(`"-%s"..parser({%s})`, f, strings.Join(p, ", ")))
		} else {
			flags = append(flags, fmt.Sprintf(`"-%s"`, f))
		}
	}

	fmt.Printf(completionScriptClinkStr,
		strings.Join(targets, ", "),
		strings.Join(flags, ",\n    "),
	)
	return
}

// handleCompletionBash returns a string to be used in bash's compgen.
// Typically, the input will look like this
//   args := []string{"flash", "-target"}
func completionBash(args []string) string {
	if len(args) == 0 {
		return fmt.Sprint(strings.Join(validCommands, " "))
	}

	cur := args[len(args)-1]
	dash := "-"

	if len(args) == 1 {
		for _, x := range validCommands {
			if x == cur {
				return ""
			}
		}
		return fmt.Sprint(strings.Join(validCommands, " "))
	} else if strings.HasPrefix(cur, "-") {
		if strings.HasPrefix(cur, "--") {
			dash = "--"
		}

		f := strings.TrimPrefix(cur, dash)
		if m, ok := flagCompleteMap[f]; ok {
			return fmt.Sprint(strings.Join(m, " "))
		} else {
			keywords := []string{}
			for _, k := range getFlagCompletion() {
				keywords = append(keywords, dash+k)
			}
			return strings.Join(keywords, " ")
		}

	} else {
		prev := args[len(args)-2]
		if strings.HasPrefix(prev, "--") {
			dash = "--"
		}

		if strings.HasPrefix(prev, "-") {
			f := strings.TrimPrefix(prev, dash)
			if m, ok := flagCompleteMap[f]; ok {
				for _, v := range m {
					if v == cur {
						return ""
					}
				}
				return fmt.Sprint(strings.Join(m, " "))
			} else {
				keywords := []string{}
				for _, k := range getFlagCompletion() {
					keywords = append(keywords, dash+k)
				}
				return strings.Join(keywords, " ")
			}
		} else {
		}
	}
	return ""
}

func getFlagCompletion() []string {
	ret := sort.StringSlice{}
	for k := range flagCompleteMap {
		ret = append(ret, k)
	}
	ret.Sort()
	return ret
}

func getTargetsFromTinygoTargets() ([]string, error) {
	buf := new(bytes.Buffer)
	cmd := exec.Command("tinygo", "targets")
	cmd.Stdout = buf
	cmd.Stderr = buf

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	targets := []string{}
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		targets = append(targets, scanner.Text())
	}

	return targets, nil
}

func getProgrammers() ([]string, error) {
	programmers := []string{}

	p, err := exec.LookPath(`openocd`)
	if err == nil {
		b := filepath.Dir(p)
		matches, err := filepath.Glob(fmt.Sprintf("%s/../share/openocd/scripts/interface/*.cfg", b))
		if err != nil {
			return nil, err
		}

		for _, m := range matches {
			programmer := strings.TrimSuffix(m, filepath.Ext(m))
			programmer = filepath.Base(programmer)
			programmers = append(programmers, programmer)
		}
	}

	if len(programmers) == 0 {
		programmers = []string{"stlink-v2", "stlink-v2-1", "cmsis-dap", "jlink", "bmp", "picoprobe"}
	}

	programmers = append(programmers, "openocd")

	return programmers, nil
}

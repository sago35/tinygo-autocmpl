package main

import (
	"fmt"
	"sort"
	"strings"
)

// flagCompleteMap specifies which flags to autocomplete.
var flagCompleteMap = map[string][]string{
	"cflags":       {},
	"dumpssa":      nil,
	"gc":           {"none", "leaking", "extalloc", "conservative"},
	"heap-size":    {},
	"ldflags":      {},
	"no-debug":     nil,
	"o":            {},
	"ocd-output":   nil,
	"opt":          {"1", "2", "s", "z"},
	"panic":        {"print", "trap"},
	"port":         {},
	"print-stacks": nil,
	"printir":      nil,
	"programmer":   {"stlink-v2", "stlink-v2-1", "cmsis-dap", "jlink"},
	"scheduler":    {"none", "tasks", "coroutines"},
	"size":         {"none", "short", "full"},
	"tags":         {},
	"target":       validTargets,
	"verifyir":     nil,
	"wasm-abi":     {},
}

// validTargets is a list of completion targets for -target. It can be overridden by arguments.
var (
	validTargets = []string{
		"arduino-mega2560", "arduino-nano", "arduino-nano33", "arduino", "atmega1284p",
		"atsame54-xpro", "atsame54p20a", "bluepill", "circuitplay-bluefruit", "circuitplay-express",
		"clue_alpha", "cortex-m-qemu", "digispark", "feather-m0", "feather-m4",
		"feather-nrf52840", "gameboy-advance", "hifive1-qemu", "hifive1b", "itsybitsy-m0",
		"itsybitsy-m4", "itsybitsy-nrf52840", "maixbit", "metro-m4-airlift", "microbit-s110v8",
		"microbit", "nintendoswitch", "nrf52840-mdk", "nucleo-f103rb", "particle-3rd-gen",
		"particle-argon", "particle-boron", "particle-xenon", "pca10031", "pca10040-s132v6",
		"pca10040", "pca10056-s140v7", "pca10056", "pinetime-devkit0", "pybadge",
		"pygamer", "pyportal", "reelboard-s140v7", "reelboard", "riscv-qemu",
		"stm32f405", "stm32f4disco-1", "stm32f4disco", "teensy36", "trinket-m0",
		"wasm", "wioterminal", "x9pro", "xiao", "esp32", "esp32-wroom-32",
		"feather-stm32f405",
	}
	validCommands = []string{"build", "run", "test", "flash", "gdb", "env", "list", "clean", "help"}
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

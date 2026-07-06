# CHANGES

0.8.0

    - Updated as of tinygo 0.41.1
    - Add fish completion support (--completion-script-fish)
    - Rename completionBash to completeArgs (shared engine used by bash/zsh/fish)
    - Add --version flag (embedded via goreleaser ldflags)
    - Continue with empty targets when tinygo is not in PATH (was fatal error)
    - Run go test in CI and update GitHub Actions to current versions
    - Install fish and its completion setup in the devcontainer

0.7.0

    - Updated as of tinygo 0.35

0.6.0

    - Updated as of tinygo 0.26
    - Add $TINYGOROOT to package completion source

0.5.0

    - Add monitor subcommand and -monitor flag

0.4.0

    - Improve error handling

0.3.0

    - Add programmers to openocd

0.2.0

    - Update flags : -p, -scheduler, wasm-abi, work
    - Fix behavior when openocd does not exist

0.1.0

    - Initial release

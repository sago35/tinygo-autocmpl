# tinygo-autocmpl

`tinygo-autocmpl` adds bash/zsh/clink completion to tinygo  

`tinygo-autocmpl` only support bash, zsh and clink at the moment, but other shells like fish should be able to do the same.  
Your contributions are welcome.  

## Description

![tinygo-autocmpl](tinygo-autocmpl.gif)

You can easily try `tinygo-autocmpl` with `GitHub Codespaces`.  

![codespace](codespace.png)

## Usage

You can enable autocompletion by setting the following to `~/.bashrc` etc.  

```
# bash
$ eval "$(tinygo-autocmpl --completion-script-bash)"

# zsh
$ eval "$(tinygo-autocmpl --completion-script-zsh)"

# clink (windows)
$ tinygo-autocmpl --completion-script-clink > %LOCALAPPDATA%\clink\tinygo.lua
```

You can customize the auto-completion of the -target flag in the following way  
This allows you to use only your own targets, for example.  

```
$ cat ~/.tinygo.targets
feather-m4
xiao

$ eval "$(tinygo-autocmpl --targets ~/.tinygo.targets --completion-script-bash)"

$ tinygo flash --target
feather-m4  xiao
```

To add wioterminal to the autocompletion candidates, do this

```
$ echo wioterminal >> ~/.tinygo.targets

$ cat ~/.tinygo.targets
feather-m4
xiao
wioterminal

$ tinygo flash --target
feather-m4   wioterminal  xiao
```

## Installation

```
go install github.com/sago35/tinygo-autocmpl@latest
```

or

download from https://github.com/sago35/tinygo-autocmpl/releases/latest

or 

download from [GitHub Actions](https://github.com/sago35/tinygo-autocmpl/actions)

### Environment

* go

I tested tinygo-autocmpl in the following environments.

* ubuntu
    * bash
    * zsh
* windows
    * bash (git for windows)
    * clink (https://mridgers.github.io/clink/)

## Notice

This project uses [goreleaser](https://goreleaser.com/) for release

## LICENSE

MIT

## Author

sago35 - <sago35@gmail.com>

# tinygo-autocmpl

`tinygo-autocmpl` adds bash completion to tinygo  

`tinygo-autocmpl` only support bash at the moment, but other shells like zsh and fish should be able to do the same.  
Your contributions are welcome.  

## Description

![tinygo-autocmpl](tinygo-autocmpl.gif)

## Usage

You can enable autocompletion by setting the following to `~/.bashrc` etc.  

```
$ eval "$(tinygo-autocmpl --completion-script-bash)"
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
go get github.com/sago35/tinygo-autocmpl
```

### Environment

* go

I tested tinygo-autocmpl in the following environments.

* ubuntu + bash
* windows + bash (git for windows)

## Notice

## Author

sago35 - <sago35@gmail.com>

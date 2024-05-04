# background-changer

This is a go cli application to change the background on windows from the command line by selecting an image from a specified repo and setting it as background

## How the CLI should work

Command should be runable via CMDLine using:

```
cbg <command> <subcommand> <optional flags>
```

## How to build 

To build this application just run the build .sh script

- You can optionallly change the version tag

## Commands

These are a list of the main commands

- get
- set
- list
- search
- config
- version

- ### get

The "get" command is use to get the file location of the current background

```
cbg get
```

- ### set

The "set" command can be used to quickly change the current background, it requires a url

```
cbg set --url "<image-url>"
```

- ### list

The "list" command displays a select list of all images in your default path and allows you to select one which will then be set as your background

```
cbg list
```

- ### config

The config is for storing the path to where you would like image files to be "get and set" from. It also contains a list of urls that you would like as quick links to search for new backgrounds


The config contains sub commands

- - show

This displays the current config values

```
cbg config show
```

- - set

This allows you to set/update the default path for images

```
cbg config set "<path-to-directory>"
```

- - addUrl

This command allows you to add a url to the url store

```
cbg cofnig addUrl "<url>"
```

- - removeUrl

This command allows you to remove a url to the url store

```
cbg cofnig removeUrl "<url>"
```
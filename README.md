# OmniSync - Sync your Tools to OmniFocus
[![Build Project](https://github.com/trevorpiltch/omnifocus-sync/actions/workflows/build.yaml/badge.svg)](https://github.com/trevorpiltch/omnifocus-sync/actions/workflows/build.yaml)

[OmniFocus](https://www.omnigroup.com/omnifocus) is a great tool that's a staple in many productive workflows. However, it lacks an easy way to sync with APIs of other tools. This project is a Go script that attemps to close that gap. </br>

Given a list of sources and projects, the program will connect to each source, parse the response, and create new issues or complete existing issues. If an item is marked complete in the API, it will be marked complete in OmniFocus. Note the this is **ONLY** a one way sync. Using the API as the single source of truth seemed safer and easier to implement for a v1.0. </br>

> Many thanks to [Mikerhodes](https://github.com/mikerhodes) for his inspiration with the [github-to-omnifocus](https://github.com/mikerhodes/github-to-omnifocus) tool. I used the tool extensively before creating this and used his code for the delta functions and OmniFocus scripts.

## Supported Versions

OmniFocus Professional 3.x  or newer </br>
[Go](https://go.dev/) v1.20 or newer

## Usage

### Config

Two configuration files are needed for this program. Both files are written in JSON and are fairly simple. Note that some example configuration files (e.g. GitHub Issues, Shortcut Stories) are included in the `examples` directory. You can use these files, but remember to place in the correct directory. For the default program, it looks in `~/.config/omnisync`, so ensure any example files or any of the files below are in that directory. </br>

#### Project

The first config file is `projects.json` where the program looks to determine which project an issue is assigned to. The format of this json is

```json
[
    {
        "URL": <URL of the project specific issues>,
        "OFName": <Name of the associated project in OmniFocus>
    }
]
```

So, for example, say I wanted to add the issues of this repository to my OmniFocus project, called OmniSync. I would configure the `projects.json` file like

```json
[
    {
        "URL": "https://api.github.com/repos/trevorpiltch/omnifocus-sync",
        "OFName": "OmniSync"
    }
]
```

#### Sources

The other config file is `sources.json`, which is where the program looks to determine the source to call for items. To write a new source, add a new item in the json array with the following fields: </br>

- `Name`: the name of the source
- `URL`: the url of the source
- `Headers`: an array of key value pairs, representing the headers to attatch to the API request
- `Queries`: a string that is attached as a query at the end of the URL
- `Response`: contains `DataField` which is a string representing the name of the top level field to look for data (usually just left blank); `Title` which is the field name to look what the name of an item is; `URL` is the link to the specific issue; `Number` is the number of the issue in the source
- `Tags`: an array of strings that represent the tags associated with this source in OmniFocus

To see an example of a source,  check out `examples/sources.json`.

### Running

To run this program, first set up the configuration by completing the previous section. Then open the command line in this directory and enter `make run`, which should build and run your program.

## Adding to OmniFocus

You can add this script as a button in OmniFocus using a few steps. First, you have to create an AppleScript. The script is simple and just does:

```applescript
on run
    tell application "Terminal"
        do script "cd <PATH_TO_PROJECT>/omnifocus-sync; make run"
    end tell
end run
```

Then in OmniFocus, open right click the toolbar and select "customize toolbar". From there, you should see an option for your script: select that and add it to your toolbar.


## Contributing

If you have an idea for OmniSync or found a bug, please open an issue with the best description you have. This work is done under the MIT license so feel free to modify, distribute, and/or monitize this program.

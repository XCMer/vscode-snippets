# vscode-snippets - Maintain your snippet library

VSCode uses a JSON format for storing snippets. However, you'd want to edit your snippets in your favorite editor (say, VSCode), and also maintain it as a git repository.

## Installation

```
go get github.com/XCMer/vscode-snippets
```

The binary should now be available under `$GOPATH/bin/vscode-snippets`.

## Configuration

Create the following configuration file in your home-directory: `$HOME/vscode-snippets.yml`

There are only two parameters that you need to set up; the location of your snippets and the location of VSCode's snippet directory. Here's how my configuration looks (I'm on OSX):

```
source_path: "/Users/<USERNAME>/mysnippets"
dest_path: "/Users/<USERNAME>/Library/Application Support/Code/User/snippets"
```

The `source_path` can be arbitrary, while the `dest_path` has to be where VSCode stores its snippets. VSCode will hot-reload any file put into this directory.

## Usage

The structure of your snippets directory should be as follows:

```
mysnippets
    - css
        - basic.css
        - reset.css
    - html
        - basic.html
        - bootstrap.html
```

These are the conventions:
1. Your snippet directory should consist of top-level folders. You can't directly put snippets here.
1. The folder level can only be one level deep. Every folder in your snippet directory should directly contain snippet files, and can have any extension that would assist with syntax highlighting in your favorite editor.
1. For every directory, vscode-snippets creates a snippet file in the target directory.

## Format of the snippet file

Your snippet file is a normal file with a YAML frontmatter. Here's how it looks:

````
---
desc: Description of your snippet
prefix: Prefix of your snippet
scope: Scope of the snippet
---
body {
    padding: 0px;
    margin: 0px;
}
````

The body of the file becomes the body of your snippet. You can put in placeholders using the same syntax as VSCode allows.

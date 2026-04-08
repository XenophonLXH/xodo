# XODO - TUI TODO List
Simple [Go](https://go.dev/) based TODO list presented as a [TUI](https://en.wikipedia.org/wiki/Text-based_user_interface)


## Libraries Used
- [Bubbletea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)
- [Bubbles](https://github.com/charmbracelet/bubbles)

## Keybindings
### List View

| Key | Description                         |
|-----|-------------------------------------|
|    a  | Add a new item                      |
|   i   | Edit an item                      |
|   d   | Mark item as done                      |
|   q  | Quit                                |


### Title View

| Key | Description                         |
|-----|-------------------------------------|
|   enter  | Save and move on                      |
|   ctrl + w  | Quit                                |
|   esc  | discard                                |


### Description View

| Key | Description                         |
|-----|-------------------------------------|
|   enter  | Save and move on                      |
|   esc  | Back to Title View                                |


### Priority View

| Key | Description                         |
|-----|-------------------------------------|
|   enter  | Save and move on                      |
|   esc  | Back to Description View                                |

## FAQ
<details>
    <summary>Where are the databases stored?</summary>

Xodo will try and store it's databases in:
```bash
$XDG_DATA_HOME/xodo/databases
```
If 
```bash
$XDG_DATA_HOME
``` 
is not set it will create a 
```bash
~/xodo/databases
``` 
directory.
</details>

<details>
    <summary>How do I use multiple lists?</summary>

Simply run: 
```bash
xodo mylist

```
</details>

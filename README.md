# console-watch

`console-watch.exe` is a simple tool that watches a single file and, if that file changes,
clears the console output and reruns a predefined command.


## Usage
To watch a file and run a command after its change, simply execute
`$ console-watch -c "cat localfile.txt" "localfile.txt"`

Each time the file `localfile.txt` changes, the current console output will be `clear`ed and after that the command `$ cat localfile.txt` will executed.

## Arguments
```
Usage: console-watch.exe [--command COMMAND] [INPUT [INPUT ...]]

Positional arguments:
  INPUT

Options:
  --command COMMAND, -c COMMAND
                         Command to execute when file(s) changes.
  --help, -h             display this help and exit
```


# remtools

move backup files to trash box, and list/empty trash box.
Note that, current implementation is only for macOS.

## Usage

### `rem`

move backup files to trash box.

```sh
rem [OPTIONS] [DIR...]
OPTIONS
    -a, --all          includes hidden file.
    -d, --dry-run      dry run mode.
    -i, --inquiry      inquiry mode.
    -r, --recursive    recursive mode.
    -v, --verbose      verbose mode.

    -h, --help         print this message and exit.
    -V, --version      print version, and exit.
```

### `lsrem`

list files in trash box.

```sh
lsrem [OPTIONS]
OPTIONS
    -a, --all            print hidden files.
    -l, --long-format    print long format.

    -h, --help           print this message.
    -V, --version        print version.
```

### `remrem`

empty trash from command line.

```sh
remrem [OPTIONS]
OPTIONS
    -i, --inquiry      inquiry mode.

    -h, --help         print this message and exit.
    -V, --version      print version, and exit.
```

## install

### Golang

```sh
$ go get github.com/tamada/remtools
```

### Homebrew

```sh
$ brew tap tamada/brew
$ brew install remtools
```

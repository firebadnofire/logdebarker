# logdebarker

`logdebarker` is a UNIX-style log sanitizer tool written in Go. It censors sensitive information from logs, config files, and other text streams based on user-defined patterns. This is useful for cleaning data before exporting logs or sharing files.

## Features

* Accepts input from STDIN or a file
* Outputs to STDOUT or a specified output file
* Uses `$HOME/.blocked_words.txt` as a pattern file for matching sensitive strings
* Enforces strict file permissions on the config file (`0600` only)
* Supports custom redaction strings (e.g. `redaction: <string>`)
* Ignores comment lines beginning with `#` in both input and config

## Usage

Pipe logs through `logdebarker`:

```sh
sudo cat /var/log/someprogram/log.log | logdebarker | tee cleaned_log.txt
```

Sanitize a file to a new output file:

```sh
sudo logdebarker /etc/nginx/conf.d/somesite.conf output.txt
```

## Configuration

`logdebarker` reads its blocklist from:

```
$HOME/.blocked_words.txt
```

Each line should be a literal string to redact. Lines starting with `#` are ignored. You may optionally define one redaction string:

```
redaction: [your-censor-string]
```

You may also use `import:` to import an entire dir's text content as a set of blocked strings.

Example `.blocked_words.txt`:

```
# block API keys
redaction: <removed>
ABC123XYZ456
supersecretpassword
import: ~/.ssh
```

## Security

To protect your sensitive word list, `logdebarker` will refuse to run if `$HOME/.blocked_words.txt` is not `chmod 600`.

```sh
chmod 600 ~/.blocked_words.txt
```

## Installation

1. Clone and build:

```sh
git clone https://codeberg.org/firebadnofire/logdebarker.git
cd logdebarker
go build -o logdebarker
```

2. Move the binary to your PATH:

```sh
sudo install -Dm755 logdebarker /usr/local/bin/logdebarker
```

OR

`go install archuser.org/logdebarker@latest`

## Final notes

`Logdebarker` is not a service, so you may simply wipe out ~/.blocked_words.txt when not in use for added security. This can be done with `truncate -s 0 ~/.blocked_words.txt`

Usage is as follows:

```
Usage:
  logdebarker [input_file] [output_file]
  logdebarker < input.txt > output.txt
  inputstream | logdebarker | outputstream
```

## License

This project is released under the GPLv3 License.

---

For bug reports, suggestions, or contributions, open an issue or submit a pull request.

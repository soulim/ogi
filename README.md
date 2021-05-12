# ogi

`ogi` (Open Graph Image) is a tool that generates social images with a geopattern background.

## Install

```ShellSession
$ go get github.com/soulim/ogi
```

## Usage

```ShellSession
$ ogi --text="Hello, world." \
      --note="https://www.example.com" \
      --width=1200 \
      --height=628 \
      --pattern="nested-squares" \
  > output.png
```

Output:

![Sample output](./docs/output.png)

NOTE: The color of the background depends on `text` option.

## Contributing

PRs accepted.

## License

MIT © Alexander Sulim


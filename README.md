# Sherry CLI

This repository includes implementation of sherry CLI.
The app allows to manage sherry state and should be used to configure sherry.

## Build

```bash
go build
```

The built file will be `./shr[.exe]`.

## Usage

```bash
shr --help
```

## Development & Testing

Configuration dir should be created on demon start if it does not exist. 
CLI will check it in `~/.sherry` by default.
If you want to use custom configuration dir, you can pass it as an argument:

```bash
shr --config "<CONFIG DIR>"
```

If you have prepared demon repository as recommended in its README, you can use it for testing CLI too.
For example, you can run demon in one terminal and use CLI in another one.
If demon repository is located on the same level as CLI, you can use this command for debugging:

```bash
go build && .\shr -c ../sherry-demon/dev-config [OTHER CLI ARGS]
```

Here we build CLI and run it with custom configuration dir.

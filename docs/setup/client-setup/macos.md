# Mac OS
If installing from the official package repository doesn't work (or requires a specific version), you can download the static binaries and place them in the following paths:

```bash
~/bin
```

## File mode
The binary `kubectl` has to be executable.

```bash
chmod +x ~/bin/kubectl
```

## Add to PATH
On Linux, you can add the `~/bin` directory to your PATH environment variable to make the `kubectl` command available to all users on the system.  
If `kubectl` is placed in a different directory, you can change the path to that directory.

```bash
export PATH=$PATH:~/bin
```

## Homebrew
If you are using [Homebrew](https://brew.sh/), you can install `kubectl` with the following command:

```bash
brew install kubectl
```

## Autocomplete

### Zsh
To enable zsh autocompletion, add the following to your `~/.zshrc` file:

```bash
source <(kubectl completion zsh)
```
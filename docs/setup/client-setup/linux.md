# 1. Linux
If installing from the official package repository doesn't work (or requires a specific version), you can download the static binaries and place them in the following paths:

```bash
~/bin
```

## 1.1 File mode
The binary `kubectl` has to be executable.

```bash
chmod +x ~/bin/kubectl
```

## 1.2 Add to PATH
On Linux, you can add the `~/bin` directory to your PATH environment variable to make the `kubectl` command available to all users on the system.  
If `kubectl` is placed in a different directory, you can change the path to that directory.

```bash
export PATH=$PATH:~/bin
```

## 1.3 Autocomplete

On most Linux distributions, you have to install the `bash-completion` package to enable autocompletion.

**Debian/Ubuntu**

```bash
sudo apt-get install bash-completion
```

**CentOS/RHEL**

```bash
sudo yum install bash-completion
```

**Fedora**

```bash
sudo dnf install bash-completion
```

### 1.3.1 Bash
To enable bash autocompletion, add the following to your `~/.bashrc` file:

```bash
source <(kubectl completion bash)
```

### 1.3.2 Zsh
To enable zsh autocompletion, add the following to your `~/.zshrc` file:

```bash
source <(kubectl completion zsh)
```

## 1.4 Verify
Verify that the `kubectl` command is available at the [Verification](setup/client-setup/verify.md) page.
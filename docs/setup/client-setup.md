# Client Setup

## Command Line

The `kubectl` command is the primary command line tool for interacting with the Kubernetes API.  
Before you can use `kubectl`, you need to install it.  

It is available for Linux, macOS, and Windows.

Offical documentation: https://kubernetes.io/docs/tasks/tools/install-kubectl/

### Linux
If installing from the official package repository doesn't work (or requires a specific version), you can download the static binaries and place them in the following paths:

```bash
~/bin
```

#### File mode
The binary `kubectl` has to be executable.

```bash
chmod +x ~/bin/kubectl
```

#### Add to PATH
On Linux, you can add the `~/bin` directory to your PATH environment variable to make the `kubectl` command available to all users on the system.  
If `kubectl` is placed in a different directory, you can change the path to that directory.

```bash
export PATH=$PATH:~/bin
```

#### Autocomplete

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

##### Bash
To enable bash autocompletion, add the following to your `~/.bashrc` file:

```bash
source <(kubectl completion bash)
```

##### Zsh
To enable zsh autocompletion, add the following to your `~/.zshrc` file:

```bash
source <(kubectl completion zsh)
```


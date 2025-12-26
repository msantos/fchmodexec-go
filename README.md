# SYNOPSIS

fchmodexec *mode* *fd* *...* -- *cmd* *...*

# DESCRIPTION

fchmodexec - fchmod(2) inherited file descriptors before exec(3)'ing
a command

`fchmodexec` runs as part of an exec chain to change the permissions of
any file descriptors inherited from the parent process before executing
a program.

# EXAMPLES

```
$ umask 077; ( ./fchmodexec 644 1 -- ls -al out) > out
$ cat out
-rw-r--r-- 1 msantos msantos 0 Mar 18 07:32 out
$ ls -al out
-rw-r--r-- 1 msantos msantos 48 Mar 18 07:32 out
```

# OPTIONS

None.

# BUILDING

```
go install go.iscode.ca/fchmodexec/cmd/fchmodexec@latest
```

To build a reproducible executable from the git repository:

```
CGO_ENABLED=0 go build -trimpath -ldflags "-w" ./cmd/fchmodexec
```

# ALTERNATIVES

- [fchmodexec](https://github.com/msantos/fchmodexec)
- [fchmodexec-rs](https://github.com/msantos/fchmodexec-rs)

# SEE ALSO

*exec(3)*, *fchmod*(2)

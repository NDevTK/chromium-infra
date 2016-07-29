# `infra/tools/protoc/`
## Building Instructions.

This instructions explain how to build `infra/tools/protoc/linux-amd64` and
`infra/tools/protoc/mac-amd64`. The steps are mostly identical.

> **Mac**. You'd need working Xcode command line tools. Then install autoconf
> and friends via your favorite package manager: `brew install autoconf automake
> libtool`

Choose a build directory. We'll use the enviornment variable `$ROOT` to
represent it.

    $ cd $ROOT

Use a release (currently using this one):

    $ curl https://github.com/google/protobuf/archive/v3.0.0.tar.gz

To build from repository version, clone the `protobuf` repository:

    $ git clone https://github.com/google/protobuf

Run the build. We include the following additional configuration flags:

- `--disable-shared`: Build a static library.
- `--prefix $ROOT/PREFIX`: The installation prefix. This allows us to install
                           locally, so we don't need root.

Build:

    $ cd protobuf
    $ ./autogen.sh
    $ ./configure --disable-shared --prefix $ROOT/PREFIX
    $ make -j24 install

This will install the generator to `$ROOT/PREFIX`. Afterwards, strip the binary.
This removes symbols and debugging information, including your username, from
the binary, reducing its size from ~40MiB to ~3MiB.

> **Mac**. Omit `-g` here.

    $ strip -g $ROOT/PREFIX/bin/protoc

We need to package the `$ROOT/PREFIX/include` directory because it includes
the protobuf standard library. However, we don't want to include all of the
C++ header files, nor do we want to include the compiled libraries. Prune the
contents of `$ROOT/PREFIX` to include only:

- The `protoc` binary.
- The standard library header files.

```
$ rm -rf $ROOT/PREFIX/lib
$ find $ROOT/PREFIX/include -type f ! -name '*.proto' -delete
```

The `protoc` utility searches for its default `include` path relative to its
binary location. First, it searches to see if the `include` path is in the same
directory as itself; if not, it looks in the parent directory. Because of the
way CIPD packages are installed, we opt to exploit the former location by
moving the `include` directory into the `bin` directory, then making the `bin`
directory the CIPD package root.

    $ mv $ROOT/PREFIX/include $ROOT/PREFIX/bin

This will result in a deployment that looks like:

- `/protoc`: The statically-linked protocol buffers compiler.
- `/include/...`: Standard `proto3` include protobufs.

Create package and deploy to CIPD server. Tag it with the Git commit of the
`protobuf` source from which it was built. Specify `<package-name>` based on
platform:

- Linux/64-bit: `infra/tools/protoc/linux-amd64`
- Mac/64-bit: `infra/tools/protoc/mac-amd64`

```
$ cipd create \
    -name <package-name> \
    -in $ROOT/PREFIX/bin \
    -install-mode copy \
    -tag "protobuf_version:<version>"
```

NOTE: If building from ToT, use a "git_commit" tag instead:

```
    ...
    -tag "git_commit:`git -C $ROOT/protobuf rev-parse HEAD`"
```

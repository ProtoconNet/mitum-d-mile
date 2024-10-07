### mitum-d-mile

*mitum-d-mile* is a d-mile contract model based on the second version of mitum(aka [mitum2](https://github.com/ProtoconNet/mitum2)).

#### Installation

Before you build `mitum-d-mile`, make sure to run `docker run` for digest api.

```sh
$ git clone https://github.com/ProtoconNet/mitum-d-mile

$ cd mitum-d-mile

$ go build -o ./mitum-d-mile
```

#### Run

```sh
$ ./mitum-d-mile init --design=<config file> <genesis file>

$ ./mitum-d-mile run <config file> --dev.allow-consensus
```

[standalong.yml](standalone.yml) is a sample of `config file`.

[genesis-design.yml](genesis-design.yml) is a sample of `genesis design file`.

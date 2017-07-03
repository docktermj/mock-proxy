# mock-proxy

Build `mock-proxy-M.m.P-I.x86_64.rpm`
and   `mock-proxy_M.m.P-I_amd64.deb`
where "M.m.P-I" is Major.minor.Patch-Iteration.

## Usage

That acts as a proxy.

It is a server of requests, and forwards them as if it's a client.

### Invocation

```console
mock-proxy socket --inbound-network <network_type> --inbound-address <inboundAddress> --outbound-network <network_type>  --outbound-address <outboundAddress>
```

## Development

### Dependencies

#### Set environment variables

```console
export GOPATH="${HOME}/go"
export PATH="${PATH}:${GOPATH}/bin:/usr/local/go/bin"
export PROJECT_DIR=${GOPATH}/src/github.com/docktermj
```

#### Download project

```console
mkdir -p ${PROJECT_DIR}
cd ${PROJECT_DIR}
git clone git@github.com:docktermj/mock-proxy.git
```

#### Download dependencies

```console
cd ${PROJECT_DIR}/mock-proxy
make dependencies
```

### Build

#### Local build

```console
cd ${PROJECT_DIR}/mock-proxy
make build-local
```

The results will be in the `${GOPATH}/bin` directory.

#### Docker build

```console
cd ${PROJECT_DIR}/mock-proxy
make build
```

The results will be in the `.../target` directory.

### Test

```console
cd ${PROJECT_DIR}/mock-proxy
make test-local
```

### Install

#### RPM-based

Example distributions: openSUSE, Fedora, CentOS, Mandrake

##### RPM Install

Example:

```console
sudo rpm -ivh mock-proxy-M.m.P-I.x86_64.rpm
```

##### RPM Update

Example: 

```console
sudo rpm -Uvh mock-proxy-M.m.P-I.x86_64.rpm
```

#### Debian

Example distributions: Ubuntu

##### Debian Install / Update

Example:

```console
sudo dpkg -i mock-proxy_M.m.P-I_amd64.deb
```

### Cleanup

```console
make clean
```

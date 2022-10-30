# RoL is Rack of Labs

RoL is an open source platform for managing bare metal servers in your lab.
We use a REST API with a CRUD model to manage configuration and interact with the system.
RoL is currently under **active development**, you can follow the current status here.

## Current status to MVP

- [x] Multi-layer architecture
- [x] Logging to database
- [x] Logs reading from API
- [x] Custom typed errors
- [x] Ethernet Switches configuration management
- [x] Ethernet Switches VLAN's and POE management (only a few switch models)
- [x] Host VLAN's and bridges management
- [x] Host network configuration saver and recover
- [x] Device templates
- [x] DHCP servers management
- [x] TFTP servers management
- [ ] Devices management
- [ ] Projects management
- [ ] iPXE provisioning

## Install Dependencies

The following steps are required to build and run RoL from source:

### Get Go 1.18.x or newest

`https://golang.org/dl/`

## How to build

`cd src && go mod tidy && go build`

## How to run

1. Add rights for network management.
   1. For RoL binary: `sudo setcap cap_net_admin+ep ./rol`
   2. For iptables run without root: `sudo setcap "cap_net_raw+ep cap_net_admin+ep" /usr/sbin/xtables-nft-multi`
   3. If you don't want to add right to iptables: you can run `./rol` as root.
2. Run RoL binary.
`./rol`
3. If all ok the last output string will be: `[GIN-debug] Listening and serving HTTP on localhost:8080`
4. Go to the [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) to read API swagger documentation.

## For developers

A typical multi-layer architecture is implemented.

### Folders structure

    .
    ├── docs                # Docs files
    │   ├── plantuml        # Struct diagrams in puml and svg formats
    ├── src                 # Source code
    │   ├── tests           # Unit tests
    │   ├── domain          # Entities
    │   ├── dtos            # DTO's is Data transfer objects
    │   ├── app             # Application logic
    │   │   ├── errors      # Custom errors implementation
    │   │   ├── interfaces  # All defined interfaces
    │   │   ├── mappers     # DTO to Entity converters
    │   │   ├── services    # Entities management logic
    │   │   ├── utils       # Utils and simple helpers
    │   │   ├── validators  # DTO's validators
    │   ├── webapi          # HTTP WEP API application
    │   │   ├── controllers # API controllers
    │   │   ├── swagger     # Swagger auto-generated docs
    │   ├── infrastructure  # Implemenatations
    └── ...

### How to update swagger docs

All description for swagger documentation is in the controllers, which are located in the src/webapi/controllers folder, changing the description of the controllers, you need to update the swagger documentation with the command:

`cd src && swag init -o webapi/swagger`

The swag utility can be installed according to the official [documentation](https://github.com/swaggo/swag#getting-started).

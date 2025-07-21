# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build Commands
- `make` or `make nginx-build` - Build the nginx-build binary
- `make build-example` - Build nginx using example configuration files
- `make check` - Run all tests
- `make fmt` - Format all Go code
- `make clean` - Remove the nginx-build binary

### Running Tests
- `go test ./...` - Run all tests in the project
- To run tests for a specific package: `go test ./builder` or `go test ./configure`

## Architecture Overview

nginx-build is a Go tool that simplifies building nginx with custom configurations and third-party modules. The codebase is organized into several key packages:

### Core Components

1. **Main Entry Point** (`nginx-build.go`): 
   - Handles command-line flags and orchestrates the build process
   - Manages parallel downloads of nginx and dependencies
   - Coordinates the configure and build phases

2. **Builder Package** (`builder/`):
   - `builder.go`: Core builder logic for different components (nginx, OpenSSL, PCRE, etc.)
   - `nginx.go`: nginx-specific build operations
   - `staticlibrary.go`: Handles static library dependencies
   - Components supported: nginx, OpenResty, freenginx, PCRE, OpenSSL, LibreSSL, zlib

3. **Configure Package** (`configure/`):
   - `configure.go`: Executes nginx configure script
   - `configureopt.go`: Manages configure options and flags
   - `generate.go`: Generates the configure script with all options
   - `normalize.go`: Normalizes configure scripts

4. **Module3rd Package** (`module3rd/`):
   - `module3rd.go`: Core third-party module management
   - `download.go`: Downloads modules from git/hg repositories
   - `load.go`: Loads module configurations from JSON
   - `provide.go`: Handles module provisioning scripts

5. **Utility Functions** (`util/`):
   - `file.go`: File operations and helpers
   - `patch.go`: Applies patches to nginx source
   - `message.go`: Error message formatting

### Key Design Patterns

1. **Parallel Downloads**: Uses goroutines and sync.WaitGroup to download nginx, dependencies, and third-party modules concurrently

2. **Builder Pattern**: Each component (nginx, PCRE, etc.) has a Builder struct that encapsulates download URLs, versions, and build logic

3. **Static Library Support**: Allows embedding PCRE, OpenSSL/LibreSSL, and zlib statically into the nginx binary

4. **Module System**: Supports both static and dynamic third-party modules via JSON configuration

5. **Work Directory Structure**: Creates organized directories under the work path:
   - `work/nginx/{version}/` for nginx builds
   - `work/openresty/{version}/` for OpenResty builds
   - `work/freenginx/{version}/` for freenginx builds

## Important Configuration Files

- `config/configure.example` - Example nginx configure script
- `config/modules.json.example` - Example third-party modules configuration
- `config/modules.json.njs` - njs module configuration example
- `config/modules.json.brotli` - Brotli module configuration example
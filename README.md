# Crow Build System

A lightweight incremental build engine written in Go. 

The project is designed as a fast and focused alternative to complex build automation tools. Its primary goal is to provide reliable task execution with strict avoidance of redundant compilation for untouched source files.

## Features

* **Incremental Execution:** Tracks task inputs and outputs using SHA-256 hashing. Tasks are skipped automatically if no source changes are detected.
* **Concurrency:** Utilizes Go's native goroutines for parallel task execution and internal operations.
* **Local Caching:** Caches dependencies and build artifacts locally to support offline workflows.
* **Declarative Configuration:** Uses a simplified configuration structure to define build steps and dependencies.

## Project Structure

* `crow/` — Core orchestration and task execution engine.
* `pkg/` — Internal packages, helper utilities, and hashing logic.
* `config/` — Configuration file parser.
* `build.crw` — Example configuration file for the build system.

## License

This project is licensed under the GNU General Public License v3 (GPL-3.0).
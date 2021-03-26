# News App Exercise - Articles Management Service

Please check out documentation for the whole system [here](https://github.com/gustavooferreira/news-app-docs).

This repository structure follows [this convention](https://github.com/golang-standards/project-layout).

---

## Tip

> If you run `make` without any targets, it will display all options available on the makefile followed by a short description.

---

# Build

To build a binary, run:

```bash
make build
```

The `api-server` binary will be placed inside the `bin/` folder.

---

# Tests

To run tests:

```bash
make test
```

To get coverage:

```bash
make coverage
```

---

# Docker

To build the docker image, run:

```bash
make build-docker
```

The docker image is named `news-app/articles-api-server`.

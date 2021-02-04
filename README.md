# bigstruct

`bigstruct` stores all your structured data in a format agnostic way. It provides out-of-the-box support for the following file formats:

- ini/toml configuration files
- yaml/json documents
- hcl/hcl+json templates

Every file is transformed into Intermediate Struct Representation (the `bigstruct/isr` package) and persisted as a list of key-value pairs in a database of your choice, either SQL or NoSQL.

The following backends are supported:

- ScyllaDB
- MySQL, sqlite, PostgreSQL (everything that GORM supports)
- etcd

## Features

### Overlays and overrides (Log-structured Merge-Tree)

### Templating

### Validation

### JSON Schema

## The `isr` format

### Codecs

### Schema

## Concepts

### Namespaces

### Indexes

## Getting started

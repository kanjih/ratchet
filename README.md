ratchet
========
ratchet is a schema migration tool for Cloud Spanner.

## Installation

Download the binary from [GitHub Releases][release] and drop it in your `$PATH`.

- [Darwin / Mac][release]
- [Linux][release]

[release]: https://github.com/hiracchy/ratchet/releases/latest

## Usage
You can use by following steps.

### 1. Initialize the schema for migration
```console
$ ratchet init -p {your-project-id} -i {spanner-instance} -d {spanner-database}

Creating migration table...
Migration table has been created!!
```

### 2. Create migration files
```console
$ ratchet new

New migration file has been created in migrations/2021-02-07_07-38-03_23229.sql
```
The above command makes a migration file for DDL.
If you want to make files for DML or Partitioned-DML, please add `--dml` or `--pdml` opition.

### 3. Run migrations
```console
$ ratchet run -p {your-project-id} -i {spanner-instance} -d {spanner-database}

Migration started.
running 2021-02-07_07-38-03_23229 ... done.
Migration completed!
```
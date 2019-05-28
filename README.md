# terraform-encrypt

A simple encrypt program to be used by terraform's external data provider or as a CLI tool.

## Usage

There are several commands which you can invoke on terraform-encrypt.

### Command: encrypt

To encrypt a file in-place (or to another file) you run:

```bash
terraform-encrypt encrypt [sourceFiles...] [flags]
```

Flags:
 - `-o`, `--output string`:     The target file location. Can only be used if a single file is passed. Specify '-' to output to stdout.
 - `-p`, `--password string`:   The vault password. This defaults to the value of environment variable `VAULT_PASSWORD`.

### Command: decrypt

To decrypt a file you run:

```bash
terraform-encrypt decrypt [sourceFiles...] [flags]
```

Flags:
 - `-c`, `--confirm-password`:  Confirm the vault password when prompting.
 - `-o`, `--output string`:     The target file location. Can only be used if a single file is passed. Specify '-' to output to stdout.
 - `-p`, `--password string`:   The vault password. This defaults to the value of environment variable `VAULT_PASSWORD`.

### Using Terraform

Create a json file:

```json
{
    "fieldA": "Value",
    "message": "I am super secret!"
}
```

Encrypt the file:

```bash
terraform-encrypt encrypt secret.json
```

Read using terraform:

```hcl
data "external" "secret" {
  program = [
    "terraform-encrypt",
    "decrypt",
    "${path.module}/path/to/encrypted/file",
    "--output",
    "-"
  ]
}

output "result" {
  value = "${data.external.secret.result.message}"
}
```

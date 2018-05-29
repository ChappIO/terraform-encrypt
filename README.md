# terraform-encrypt
A simple encrypt program to be used by terraform's external data provider

## Usage

There are several commands which you can invoke on terraform-encrypt.

### Using Terraform

````hcl

data "external" "secret" {
  program = [
    "terraform-encrypt"
  ]
  
  query = {
    vault_key = "My Secret Key... 123"
    src_file = "${path.root}/path/to/file.json"
  }
}

output "result" {
  value = "${data.external.ansible.result.message}"
}

````

### Command: terraform-encrypt encrypt

To encypt a file in-place (or to another file) you run:

```bash
terraform-encrypt encrypt <source-file>
```

Options:

- `-o <target-file>`: Store the encypted version of the source file in a different place

### Command: terraform-encrypt decrypt

### Command: terraform-encrypt job
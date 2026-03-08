# terraform-provider-namingservice (Public)

Standalone Terraform provider project for public distribution.

This version keeps the same naming inputs as your internal provider, and requires `api_url` as a resource input so any user can point to their own naming service endpoint.

## Resource

Resource name: `namingservice_name`

Inputs:
- `api_url` (required): Full endpoint URL, example `https://api.example.com/generate-name`
- `app` (required)
- `env` (required)
- `region` (required)
- `csp` (required)
- `resource_type` (required)

Outputs:
- `generated_name` (computed)
- `id` (computed)

## Local Build

From `public-provider/terraform-provider-namingservice`:

```powershell
go mod tidy
go build -o terraform-provider-namingservice.exe
```

## Example Usage

Use a local source while developing:

```hcl
terraform {
  required_providers {
    namingservice = {
      source  = "local/namingservice"
      version = "0.1.0"
    }
  }
}

provider "namingservice" {}

resource "namingservice_name" "rg" {
  api_url       = "http://localhost:8000/generate-name"
  app           = "pay"
  env           = "dev"
  region        = "weu"
  csp           = "az"
  resource_type = "rg"
}

output "generated_rg_name" {
  value = namingservice_name.rg.generated_name
}
```

## Prepare for Public Publishing

1. Move this provider to its own public GitHub repository, for example:
   - `github.com/<your-org>/terraform-provider-namingservice`
2. Update `go.mod` module path to match that repo.
3. Tag releases with semantic versions (`v0.1.0`, `v0.1.1`, etc.).
4. Build release binaries for all target OS/arch.
5. Publish checksums and signed artifacts.

## Public Terraform Registry Notes

To publish on Terraform Registry, use source format:
- `<namespace>/namingservice`

Example after publish:

```hcl
terraform {
  required_providers {
    namingservice = {
      source  = "<namespace>/namingservice"
      version = "0.1.0"
    }
  }
}
```

Registry requirements typically include:
- Public repository named `terraform-provider-<type>`
- Git tags for releases
- Release artifacts for supported platforms
- Provider docs and examples

## Next Improvements

- Add resource validation for URL and short-code formats.
- Add retries and configurable timeout.
- Add auth support (token/header).
- Add acceptance tests.

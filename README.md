# terraform-provider-qwen

Terraform provider for Alibaba Cloud Qwen / DashScope admin API — workspaces and API keys.

## Usage

```hcl
terraform {
  required_providers {
    qwen = {
      source  = "stark256-spec/qwen"
      version = "~> 1.0"
    }
  }
}

provider "qwen" {
  api_key = var.dashscope_api_key
}

resource "qwen_workspace" "prod" {
  name = "production"
}

resource "qwen_api_key" "app" {
  name = "app-server"
}
```

## Authentication

Set your API key via the `api_key` argument or the environment variable shown in the provider schema.

## Resources

| Resource | Description |
|----------|-------------|
| `qwen_workspace` / `qwen_project` / `qwen_team` | Isolated environment |
| `qwen_api_key` | API key scoped to a workspace/project |

## License

Apache 2.0

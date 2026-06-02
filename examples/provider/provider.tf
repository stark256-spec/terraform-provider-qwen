terraform {
  required_providers {
    qwen = {
      source  = "stark256-spec/qwen"
      version = "~> 1.0"
    }
  }
}

provider "qwen" {
  api_key = var.dashscope_api_key   # or DASHSCOPE_API_KEY env var
}

resource "qwen_workspace" "prod" {
  name = "production"
}

resource "qwen_api_key" "app" {
  name = "app-server"
}

output "app_key" {
  value     = qwen_api_key.app.secret_key
  sensitive = true
}
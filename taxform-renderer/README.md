# taxform-renderer

HTTP service that fills PDF form templates with provided data. Instrumented with OpenTelemetry.

## Usage

```bash
go run .
```

### Fill a form

```bash
curl -X POST http://localhost:8080/render \
  -H 'Content-Type: application/json' \
  -d '{
    "template": "default",
    "fields": {
      "nome completo": "João Silva",
      "data de nascimento": "01/01/1990",
      "nacionalidade": "Brasileira",
      "nome do pai": "Carlos Silva",
      "nome da mãe": "Maria Silva",
      "endereço completo rua número município  local país": "Rua Example 123, São Paulo, Brasil",
      "Local e Data": "São Paulo, 19/03/2026"
    },
    "radioButtonGroups": {
      "Condição": "residente"
    }
  }' -o filled.pdf
```

### Request format

| Field | Type | Description |
|---|---|---|
| `template` | string | PDF template name (without `.pdf`), defaults to `default` |
| `fields` | map | Text field name/value pairs |
| `dateFields` | map | Date field name/value pairs |
| `radioButtonGroups` | map | Radio button group name/selected option pairs |

Templates are loaded from the `templates/` directory (override with `TEMPLATE_DIR` env var).

### List available fields in a template

```bash
go run github.com/pdfcpu/pdfcpu/cmd/pdfcpu@latest form list templates/default.pdf
```

## Configuration

| Env var | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `TEMPLATE_DIR` | `templates` | Path to PDF templates |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `localhost:4317` | OTel collector gRPC endpoint |

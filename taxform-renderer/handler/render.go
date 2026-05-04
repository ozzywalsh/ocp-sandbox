package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
)

var tracer = otel.Tracer("taxform-renderer/handler")
var meter = otel.Meter("taxform-renderer/handler")
var renderCount, _ = meter.Int64Counter("render.count",
	metric.WithDescription("Total number of tax form renders"),
)

type RenderRequest struct {
	Template          string            `json:"template"`
	Fields            map[string]string `json:"fields"`
	DateFields        map[string]string `json:"dateFields,omitempty"`
	RadioButtonGroups map[string]string `json:"radioButtonGroups,omitempty"`
}

func Render() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "render")
		defer span.End()

		var req RenderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.WarnContext(ctx, "invalid request body", "error", err)
			span.SetStatus(codes.Error, "invalid request body")
			renderCount.Add(ctx, 1,
				metric.WithAttributes(attribute.String("status", "error"), attribute.String("template", "")),
			)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.Template == "" {
			req.Template = "default"
		}
		span.SetAttributes(attribute.String("template", req.Template))

		templateDir := os.Getenv("TEMPLATE_DIR")
		if templateDir == "" {
			templateDir = "templates"
		}
		templatePath := filepath.Join(templateDir, req.Template+".pdf")

		pdf, err := fillPDF(ctx, templatePath, req)
		if err != nil {
			slog.ErrorContext(ctx, "failed to render pdf", "error", err, "template", req.Template)
			span.SetStatus(codes.Error, err.Error())
			renderCount.Add(ctx, 1,
				metric.WithAttributes(attribute.String("status", "error"), attribute.String("template", req.Template)),
			)
			http.Error(w, fmt.Sprintf("failed to render pdf: %v", err), http.StatusInternalServerError)
			return
		}

		renderCount.Add(ctx, 1,
			metric.WithAttributes(attribute.String("status", "success"), attribute.String("template", req.Template)),
		)
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-filled.pdf", req.Template))
		w.Write(pdf)
	})
}

type formField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func fillPDF(ctx context.Context, templatePath string, req RenderRequest) ([]byte, error) {
	_, span := tracer.Start(ctx, "fillPDF")
	defer span.End()

	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	form := map[string]any{}

	if len(req.Fields) > 0 {
		tf := make([]formField, 0, len(req.Fields))
		for k, v := range req.Fields {
			tf = append(tf, formField{Name: k, Value: v})
		}
		form["textfield"] = tf
	}

	if len(req.DateFields) > 0 {
		df := make([]formField, 0, len(req.DateFields))
		for k, v := range req.DateFields {
			df = append(df, formField{Name: k, Value: v})
		}
		form["datefield"] = df
	}

	if len(req.RadioButtonGroups) > 0 {
		rbg := make([]formField, 0, len(req.RadioButtonGroups))
		for k, v := range req.RadioButtonGroups {
			rbg = append(rbg, formField{Name: k, Value: v})
		}
		form["radiobuttongroup"] = rbg
	}

	formJSON := map[string]any{
		"forms": []map[string]any{form},
	}

	jsonData, err := json.Marshal(formJSON)
	if err != nil {
		return nil, fmt.Errorf("marshal form data: %w", err)
	}

	in := bytes.NewReader(templateData)
	rd := bytes.NewReader(jsonData)
	var out bytes.Buffer

	if err := api.FillForm(in, rd, &out, nil); err != nil {
		return nil, fmt.Errorf("fill form: %w", err)
	}

	totalFields := len(req.Fields) + len(req.DateFields) + len(req.RadioButtonGroups)
	span.SetAttributes(attribute.Int("fields.count", totalFields))
	return out.Bytes(), nil
}

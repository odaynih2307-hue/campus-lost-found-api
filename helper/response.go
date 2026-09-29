package helper

import (
	"campus-lost-found-api/app/model"
	"context"
	"encoding/csv"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type errorEnvelope struct {
	Success   bool              `json:"success"`
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"request_id"`
}

func Success(c *fiber.Ctx, status int, message string, data, meta any) error {
	return c.Status(status).JSON(envelope{
		Success: status >= 200 && status < 300,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	ae, ok := err.(*APIError)
	if !ok {
		ae = Internal(err)
	}

	if ae.Cause != nil {
		_ = ae.Cause
	}

	requestID, _ := c.Locals("request_id").(string)

	return c.Status(ae.Status).JSON(errorEnvelope{
		Success:   false,
		Code:      ae.Code,
		Message:   ae.Message,
		Fields:    ae.Fields,
		RequestID: requestID,
	})
}

func Validate(v *validator.Validate, input any) error {
	if err := v.Struct(input); err != nil {
		fields := map[string]string{}

		for _, fe := range err.(validator.ValidationErrors) {
			fields[strings.ToLower(fe.Field())] = fe.Tag()
		}

		return Validation(fields)
	}

	return nil
}

func ParseID(c *fiber.Ctx) (int64, error) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)

	if err != nil || id <= 0 {
		return 0, BadID()
	}

	return id, nil
}

func RequestContext(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), 5*time.Second)
}

func Negotiate(c *fiber.Ctx, offered ...string) (string, error) {
	accept := strings.TrimSpace(c.Get(fiber.HeaderAccept))

	if accept == "" || accept == "*/*" {
		return offered[0], nil
	}

	chosen := c.Accepts(offered...)

	if chosen == "" {
		return "", NotAcceptable("format yang diminta tidak tersedia")
	}

	return chosen, nil
}

func WriteItemsCSV(c *fiber.Ctx, items []model.Item) error {
	c.Set(fiber.HeaderContentType, "text/csv; charset=utf-8")
	c.Set(fiber.HeaderContentDisposition, `attachment; filename="items.csv"`)

	var b strings.Builder

	w := csv.NewWriter(&b)

	_ = w.Write([]string{
		"id",
		"title",
		"description",
		"category",
		"location",
		"status",
		"reported_by",
		"created_at",
	})

	for _, i := range items {
		_ = w.Write([]string{
			strconv.FormatInt(i.ID, 10),
			i.Title,
			i.Description,
			i.Category,
			i.Location,
			i.Status,
			strconv.FormatInt(i.ReportedBy, 10),
			i.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return Internal(err)
	}

	return c.SendString(b.String())
}

var _ = json.Valid
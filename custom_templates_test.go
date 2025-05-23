package maildoor

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCustomTemplates(t *testing.T) {
	t.Run("custom layout template", func(t *testing.T) {
		customLayout := `<!DOCTYPE html>
<html>
<head>
    <title>Custom Layout - {{.ProductName}}</title>
</head>
<body>
    <div class="custom-wrapper">
        {{block "yield" .}}{{end}}
    </div>
</body>
</html>`

		handler := New(
			CustomLayoutTemplate(customLayout),
			ProductName("Test App"),
		)

		req := httptest.NewRequest("GET", "/login", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		body := w.Body.String()
		if !strings.Contains(body, "Custom Layout - Test App") {
			t.Error("Custom layout template not used")
		}
		if !strings.Contains(body, `<div class="custom-wrapper">`) {
			t.Error("Custom layout wrapper not found")
		}
	})

	t.Run("custom login template", func(t *testing.T) {
		customLogin := `{{block "title" .}}Custom Login - {{.ProductName}}{{end}}

{{define "yield"}}
<div class="custom-login-form">
    <h1>Custom Login Page</h1>
    <p>Welcome to {{.ProductName}}</p>
    <form action="{{prefixedPath "/email"}}" method="POST">
        <input type="email" name="email" placeholder="Enter your email" required>
        <button type="submit">Custom Login Button</button>
    </form>
    {{if ne .Error ""}}
        <div class="custom-error">{{.Error}}</div>
    {{end}}
</div>
{{end}}`

		handler := New(
			CustomLoginTemplate(customLogin),
			ProductName("Custom App"),
		)

		req := httptest.NewRequest("GET", "/login", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		body := w.Body.String()
		if !strings.Contains(body, "Custom Login Page") {
			t.Error("Custom login template not used")
		}
		if !strings.Contains(body, "Welcome to Custom App") {
			t.Error("Custom login template data not rendered")
		}
		if !strings.Contains(body, "Custom Login Button") {
			t.Error("Custom login form not found")
		}
	})

	t.Run("both custom layout and login templates", func(t *testing.T) {
		customLayout := `<!DOCTYPE html>
<html>
<head>
    <title>{{block "title" .}}Default{{end}}</title>
    <style>.custom-style { color: blue; }</style>
</head>
<body class="custom-body">
    {{block "yield" .}}{{end}}
</body>
</html>`

		customLogin := `{{block "title" .}}Both Custom - {{.ProductName}}{{end}}

{{define "yield"}}
<div class="custom-style">
    <h2>Both Templates Custom</h2>
    <form action="{{prefixedPath "/email"}}" method="POST">
        <input type="email" name="email" required>
        <button type="submit">Submit</button>
    </form>
</div>
{{end}}`

		handler := New(
			CustomLayoutTemplate(customLayout),
			CustomLoginTemplate(customLogin),
			ProductName("Dual Custom"),
		)

		req := httptest.NewRequest("GET", "/login", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		body := w.Body.String()
		if !strings.Contains(body, "Both Custom - Dual Custom") {
			t.Error("Custom title not rendered")
		}
		if !strings.Contains(body, `<body class="custom-body">`) {
			t.Error("Custom layout body not found")
		}
		if !strings.Contains(body, "Both Templates Custom") {
			t.Error("Custom login content not found")
		}
		if !strings.Contains(body, `.custom-style { color: blue; }`) {
			t.Error("Custom layout styles not found")
		}
	})

	t.Run("custom template with error handling", func(t *testing.T) {
		customLogin := `{{define "yield"}}
<div class="error-test">
    {{if ne .Error ""}}
        <span class="custom-error-message">Error: {{.Error}}</span>
    {{end}}
    <form action="{{prefixedPath "/email"}}" method="POST">
        <input type="email" name="email" required>
        <button type="submit">Test Submit</button>
    </form>
</div>
{{end}}`

		handler := New(
			CustomLoginTemplate(customLogin),
			EmailValidator(func(email string) error {
				if email == "invalid@example.com" {
					return io.ErrUnexpectedEOF
				}
				return nil
			}),
		)

		// Test with invalid email to trigger error
		req := httptest.NewRequest("POST", "/email", strings.NewReader("email=invalid@example.com"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		body := w.Body.String()
		if !strings.Contains(body, "custom-error-message") {
			t.Error("Custom error template not used")
		}
		if !strings.Contains(body, "unexpected EOF") {
			t.Error("Error message not displayed in custom template")
		}
	})

	t.Run("custom template with prefix", func(t *testing.T) {
		customLogin := `{{define "yield"}}
<form action="{{prefixedPath "/email"}}" method="POST" class="prefixed-form">
    <input type="email" name="email" required>
    <button type="submit">Submit</button>
</form>
{{end}}`

		handler := New(
			CustomLoginTemplate(customLogin),
			Prefix("/auth"),
		)

		req := httptest.NewRequest("GET", "/auth/login", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		body := w.Body.String()
		if !strings.Contains(body, `action="/auth/email"`) {
			t.Error("Prefixed path not correctly generated in custom template")
		}
		if !strings.Contains(body, "prefixed-form") {
			t.Error("Custom template not rendered")
		}
	})

	t.Run("fallback to embedded templates when custom not provided", func(t *testing.T) {
		handler := New(ProductName("Fallback Test"))

		req := httptest.NewRequest("GET", "/login", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		body := w.Body.String()
		if !strings.Contains(body, "Sign in to your account") {
			t.Error("Default embedded template not used as fallback")
		}
		if !strings.Contains(body, "Send me a login code") {
			t.Error("Default login form not found")
		}
	})

	t.Run("mixed custom and embedded templates", func(t *testing.T) {
		// Only provide custom layout, login should use embedded
		customLayout := `<!DOCTYPE html>
<html>
<head><title>Mixed Test</title></head>
<body class="mixed-layout">
    {{block "yield" .}}{{end}}
</body>
</html>`

		handler := New(
			CustomLayoutTemplate(customLayout),
			ProductName("Mixed Test"),
		)

		req := httptest.NewRequest("GET", "/login", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		body := w.Body.String()
		if !strings.Contains(body, "Mixed Test") {
			t.Error("Custom layout not used")
		}
		if !strings.Contains(body, "mixed-layout") {
			t.Error("Custom layout class not found")
		}
		if !strings.Contains(body, "Sign in to your account") {
			t.Error("Embedded login template not used")
		}
	})

	t.Run("invalid custom template syntax", func(t *testing.T) {
		invalidTemplate := `{{define "yield"}}
<div>
    {{.InvalidSyntax unclosed
</div>
{{end}}`

		handler := New(
			CustomLoginTemplate(invalidTemplate),
		)

		req := httptest.NewRequest("GET", "/login", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		// Should return 500 error due to template parsing error
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected 500 error for invalid template, got %d", w.Code)
		}
	})
}
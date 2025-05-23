# Custom Login Page Templates

Maildoor supports custom login page templates, allowing you to completely customize the appearance and styling of your login pages while maintaining all the authentication functionality.

## Overview

The custom template system allows you to:
- Replace the default layout template with your own HTML structure
- Replace the default login form with custom styling and content
- Mix custom and default templates (e.g., custom layout with default login form)
- Maintain full compatibility with all Maildoor features (prefixed paths, error handling, etc.)

## Quick Start

```go
package main

import (
    "github.com/wawandco/maildoor"
    "net/http"
)

func main() {
    customLogin := `{{define "yield"}}
    <div class="my-custom-login">
        <h1>Welcome to {{.ProductName}}</h1>
        <form action="{{prefixedPath "/email"}}" method="POST">
            <input type="email" name="email" required>
            <button type="submit">Login</button>
        </form>
        {{if ne .Error ""}}
            <div class="error">{{.Error}}</div>
        {{end}}
    </div>
    {{end}}`

    handler := maildoor.New(
        maildoor.CustomLoginTemplate(customLogin),
        maildoor.ProductName("MyApp"),
    )

    http.ListenAndServe(":8080", handler)
}
```

## Available Options

### CustomLayoutTemplate(template string)

Replace the entire HTML layout structure:

```go
customLayout := `<!DOCTYPE html>
<html>
<head>
    <title>{{block "title" .}}{{.ProductName}}{{end}}</title>
    <style>
        body { background: #f0f0f0; font-family: Arial, sans-serif; }
        .container { max-width: 400px; margin: 50px auto; }
    </style>
</head>
<body>
    <div class="container">
        {{block "yield" .}}{{end}}
    </div>
</body>
</html>`

handler := maildoor.New(
    maildoor.CustomLayoutTemplate(customLayout),
)
```

### CustomLoginTemplate(template string)

Replace just the login form content:

```go
customLogin := `{{block "title" .}}Sign In - {{.ProductName}}{{end}}

{{define "yield"}}
<div class="login-form">
    <img src="{{.Logo}}" alt="Logo">
    <h2>Please sign in</h2>
    
    {{if ne .Error ""}}
        <div class="alert alert-error">{{.Error}}</div>
    {{end}}
    
    <form action="{{prefixedPath "/email"}}" method="POST">
        <div class="form-group">
            <label for="email">Email:</label>
            <input type="email" id="email" name="email" required>
        </div>
        <button type="submit" class="btn btn-primary">
            Send Login Code
        </button>
    </form>
</div>
{{end}}`

handler := maildoor.New(
    maildoor.CustomLoginTemplate(customLogin),
)
```

## Template Data

Your custom templates have access to the following data:

| Field | Type | Description |
|-------|------|-------------|
| `.Logo` | string | Logo URL (set via `Logo()` option) |
| `.Icon` | string | Icon URL (set via `Icon()` option) |
| `.ProductName` | string | Product name (set via `ProductName()` option) |
| `.Email` | string | User's email (available when showing errors) |
| `.Error` | string | Error message (empty if no error) |

## Template Functions

| Function | Description | Example |
|----------|-------------|---------|
| `prefixedPath` | Generate URLs with the correct prefix | `{{prefixedPath "/email"}}` |

## Required Template Blocks

### Layout Template Requirements

Your custom layout template must include:

1. **`{{block "title" .}}Default Title{{end}}`** - Page title block
2. **`{{block "yield" .}}{{end}}`** - Content area where login form appears

### Login Template Requirements

Your custom login template should include:

1. **`{{define "yield"}}`** - Main content definition
2. **Form with correct action** - `<form action="{{prefixedPath "/email"}}" method="POST">`
3. **Email input** - `<input type="email" name="email" required>`
4. **Error display** - Show `{{.Error}}` when not empty

## Examples

### Example 1: Modern Dark Theme

```go
darkLayout := `<!DOCTYPE html>
<html>
<head>
    <title>{{block "title" .}}{{.ProductName}}{{end}}</title>
    <style>
        body {
            background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
            color: white;
            font-family: 'Segoe UI', sans-serif;
            margin: 0;
            padding: 0;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .card {
            background: rgba(255,255,255,0.1);
            backdrop-filter: blur(10px);
            border-radius: 15px;
            padding: 40px;
            box-shadow: 0 8px 32px rgba(0,0,0,0.3);
        }
    </style>
</head>
<body>
    <div class="card">
        {{block "yield" .}}{{end}}
    </div>
</body>
</html>`

darkLogin := `{{block "title" .}}{{.ProductName}} - Secure Login{{end}}

{{define "yield"}}
<div style="text-align: center;">
    <h1 style="margin-bottom: 30px;">{{.ProductName}}</h1>
    
    {{if ne .Error ""}}
    <div style="background: rgba(255,0,0,0.2); padding: 10px; border-radius: 5px; margin-bottom: 20px;">
        {{.Error}}
    </div>
    {{end}}
    
    <form action="{{prefixedPath "/email"}}" method="POST">
        <input 
            type="email" 
            name="email" 
            placeholder="Enter your email"
            required
            style="width: 100%; padding: 15px; margin-bottom: 20px; border: none; border-radius: 8px; background: rgba(255,255,255,0.2); color: white; font-size: 16px;"
        >
        <button 
            type="submit"
            style="width: 100%; padding: 15px; border: none; border-radius: 8px; background: #4CAF50; color: white; font-size: 16px; cursor: pointer;"
        >
            Access My Account
        </button>
    </form>
</div>
{{end}}`

handler := maildoor.New(
    maildoor.CustomLayoutTemplate(darkLayout),
    maildoor.CustomLoginTemplate(darkLogin),
    maildoor.ProductName("SecureApp"),
)
```

### Example 2: Corporate Style

```go
corporateLogin := `{{define "yield"}}
<div style="max-width: 350px; margin: 0 auto; padding: 40px 20px;">
    <div style="text-align: center; margin-bottom: 40px;">
        <img src="{{.Logo}}" alt="{{.ProductName}}" style="height: 60px; margin-bottom: 20px;">
        <h1 style="color: #333; font-size: 28px; margin: 0;">Employee Portal</h1>
        <p style="color: #666; margin: 10px 0 0;">{{.ProductName}} Secure Access</p>
    </div>

    {{if ne .Error ""}}
    <div style="border-left: 4px solid #e74c3c; background: #fdf2f2; padding: 12px 16px; margin-bottom: 20px;">
        <strong style="color: #c0392b;">Authentication Error:</strong><br>
        <span style="color: #e74c3c;">{{.Error}}</span>
    </div>
    {{end}}

    <form action="{{prefixedPath "/email"}}" method="POST">
        <div style="margin-bottom: 25px;">
            <label style="display: block; margin-bottom: 8px; color: #333; font-weight: 600;">
                Corporate Email Address
            </label>
            <input 
                type="email" 
                name="email" 
                required
                style="width: 100%; padding: 12px; border: 2px solid #ddd; border-radius: 4px; font-size: 16px; box-sizing: border-box;"
                placeholder="firstname.lastname@company.com"
            >
        </div>
        
        <button 
            type="submit"
            style="width: 100%; background: #3498db; color: white; border: none; padding: 14px; border-radius: 4px; font-size: 16px; font-weight: 600; cursor: pointer;"
        >
            Request Access Code
        </button>
    </form>
    
    <div style="text-align: center; margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee;">
        <p style="color: #999; font-size: 12px; margin: 0;">
            Protected by {{.ProductName}} Enterprise Security
        </p>
    </div>
</div>
{{end}}`

handler := maildoor.New(
    maildoor.CustomLoginTemplate(corporateLogin),
    maildoor.ProductName("CorpSystem"),
    maildoor.Logo("https://example.com/corporate-logo.png"),
)
```

### Example 3: Mixed Templates (Custom Layout + Default Form)

```go
customLayout := `<!DOCTYPE html>
<html>
<head>
    <title>{{block "title" .}}{{.ProductName}} Portal{{end}}</title>
    <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
</head>
<body class="bg-gradient-to-r from-blue-400 to-purple-500 min-h-screen flex items-center justify-center">
    <div class="bg-white rounded-lg shadow-xl p-8 max-w-md w-full">
        {{block "yield" .}}{{end}}
    </div>
</body>
</html>`

// Using only custom layout - login form will use Maildoor's default
handler := maildoor.New(
    maildoor.CustomLayoutTemplate(customLayout),
    maildoor.ProductName("MyPortal"),
)
```

## Best Practices

### 1. Maintain Accessibility
```go
// Good: Include proper labels and ARIA attributes
`<label for="email">Email Address</label>
<input type="email" id="email" name="email" required aria-describedby="email-help">
<div id="email-help">We'll send a login code to this address</div>`
```

### 2. Mobile-Responsive Design
```go
// Include viewport meta tag in custom layout
`<meta name="viewport" content="width=device-width, initial-scale=1.0">`

// Use responsive CSS
`<style>
@media (max-width: 480px) {
    .login-container { padding: 20px 15px; }
}
</style>`
```

### 3. Handle Errors Gracefully
```go
// Always check for errors and provide clear feedback
`{{if ne .Error ""}}
<div class="error-message" role="alert">
    <strong>Unable to proceed:</strong> {{.Error}}
</div>
{{end}}`
```

### 4. Security Considerations
```go
// Always use the prefixedPath function for form actions
`<form action="{{prefixedPath "/email"}}" method="POST">`

// Include CSRF protection if needed
`<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">`
```

## Fallback Behavior

- If no custom templates are provided, Maildoor uses its built-in templates
- If only a custom layout is provided, the default login form is used within your layout
- If only a custom login template is provided, it's used within the default layout
- Template parsing errors will result in HTTP 500 responses

## Migration from Default Templates

To migrate from default templates:

1. **Start with layout**: Copy the default layout and modify styling
2. **Customize login form**: Replace form HTML while keeping the same form structure
3. **Test thoroughly**: Ensure error states and form submission work correctly
4. **Validate across devices**: Test responsive behavior on different screen sizes

## Advanced Usage

### Template Inheritance
```go
// You can define reusable template parts
baseStyle := `<style>
.btn { padding: 12px 24px; border: none; border-radius: 6px; cursor: pointer; }
.btn-primary { background: #007cba; color: white; }
.form-input { width: 100%; padding: 10px; border: 1px solid #ddd; border-radius: 4px; }
</style>`

customLayout := `<!DOCTYPE html>
<html>
<head>
    <title>{{block "title" .}}{{.ProductName}}{{end}}</title>
    ` + baseStyle + `
</head>
<body>
    {{block "yield" .}}{{end}}
</body>
</html>`
```

### Dynamic Styling Based on Data
```go
customLogin := `{{define "yield"}}
<div class="login-form">
    <h1 style="color: {{if eq .ProductName "DarkApp"}}#333{{else}}#007cba{{end}};">
        {{.ProductName}}
    </h1>
    <!-- rest of template -->
</div>
{{end}}`
```

This feature provides complete flexibility while maintaining the security and functionality of Maildoor's authentication system.
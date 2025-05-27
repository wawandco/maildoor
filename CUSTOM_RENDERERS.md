# Custom Renderers Feature

This document describes the custom renderer functionality added to maildoor that allows users to completely customize the appearance of login and code entry pages.

## Overview

Maildoor now supports custom renderer functions that receive the data needed to render authentication pages and return HTML strings. This allows complete control over the visual design while maintaining maildoor's authentication logic.

## Available Renderers

### LoginRenderer

Customizes the login page where users enter their email address.

```go
maildoor.LoginRenderer(func(data maildoor.Attempt) (string, error) {
    // Return custom HTML for login page
    return htmlString, nil
})
```

### CodeRenderer  

Customizes the code entry page where users enter their verification code.

```go
maildoor.CodeRenderer(func(data maildoor.Attempt) (string, error) {
    // Return custom HTML for code page
    return htmlString, nil
})
```

## Data Structure

Both renderers receive a `maildoor.Attempt` struct containing:

- `Logo` (string) - URL of the logo image
- `Icon` (string) - URL of the icon image  
- `ProductName` (string) - Name of your product
- `Email` (string) - User's email address (available in code renderer)
- `Error` (string) - Error message if validation failed
- `Code` (string) - Verification code (context-dependent)

## Form Requirements

Custom renderers must include correct form actions and fields:

**Login Form:**
- Action: `/email` (or prefixed path)
- Method: POST
- Required field: `email`

**Code Form:**
- Action: `/code` (or prefixed path)  
- Method: POST
- Required fields: `email` (hidden), `code`

## Usage Example

```go
handler := maildoor.New(
    maildoor.ProductName("MyApp"),
    maildoor.LoginRenderer(customLoginRenderer),
    maildoor.CodeRenderer(customCodeRenderer),
    // other options...
)
```

## Fallback Behavior

- If no custom renderer is provided, default templates are used
- Renderers can be mixed (e.g., custom login, default code page)
- Renderer errors result in HTTP 500 responses

## Best Practices

1. Always check for `data.Error` and display error messages
2. Include proper form validation and accessibility features
3. Handle both success and error states appropriately
4. Return complete HTML documents with proper DOCTYPE
5. Consider mobile responsiveness in custom designs

## Testing

The feature includes comprehensive tests covering:
- Custom renderer functionality
- Error handling
- Prefix path compatibility
- Fallback to default templates
- Integration with existing maildoor features

See `examples/custom_renderers.go` for a complete implementation example.
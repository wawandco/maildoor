# Custom Renderers Example

This example demonstrates how to use custom renderer functions with maildoor to create your own login and code entry pages.

## Overview

Maildoor allows you to customize the appearance of the login and code entry screens by providing custom renderer functions. These functions receive the data needed to render the page and return an HTML string.

## Usage

```go
package main

import (
    "github.com/wawandco/maildoor"
)

func main() {
    handler := maildoor.New(
        maildoor.ProductName("MyApp"),
        maildoor.Logo("https://example.com/logo.png"),
        maildoor.LoginRenderer(customLoginRenderer),
        maildoor.CodeRenderer(customCodeRenderer),
        maildoor.EmailSender(yourEmailSender),
    )
    
    // Use handler in your HTTP server
}
```

## Renderer Functions

### Login Renderer

The login renderer function receives an `maildoor.Attempt` struct with the following fields:

- `Logo` - URL of the logo image
- `Icon` - URL of the icon image  
- `ProductName` - Name of your product
- `Error` - Error message (if any)

```go
func customLoginRenderer(data maildoor.Attempt) (string, error) {
    // Return your custom HTML here
    html := `<html>...</html>`
    return html, nil
}
```

### Code Renderer

The code renderer function receives an `maildoor.Attempt` struct with the same fields plus:

- `Email` - The email address where the code was sent
- `Code` - The verification code (only available in certain contexts)

```go
func customCodeRenderer(data maildoor.Attempt) (string, error) {
    // Return your custom HTML here
    html := `<html>...</html>`
    return html, nil
}
```

## Important Considerations

1. **Form Actions**: Make sure your forms submit to the correct endpoints:
   - Login form should POST to `/email`
   - Code form should POST to `/code`

2. **Required Fields**: 
   - Login form needs an `email` field
   - Code form needs `email` (hidden) and `code` fields

3. **Error Handling**: Check `data.Error` and display error messages appropriately

4. **Content Type**: The renderers should return complete HTML documents. Maildoor will set the `Content-Type` header to `text/html`.

## Running the Example

```bash
go run custom_renderers.go
```

Then visit `http://localhost:8080/login` to see the custom renderers in action.

## Fallback Behavior

If you don't provide custom renderers, maildoor will use its default templates. You can mix and match - for example, only customize the login page while keeping the default code entry page.
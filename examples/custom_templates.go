package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wawandco/maildoor"
)

func main() {
	// Example 1: Custom layout template with dark theme
	customLayout := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{block "title" .}}{{.ProductName}} - Login{{end}}</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            margin: 0;
            padding: 0;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .container {
            background: white;
            border-radius: 12px;
            box-shadow: 0 10px 25px rgba(0,0,0,0.1);
            overflow: hidden;
            max-width: 400px;
            width: 100%;
        }
    </style>
</head>
<body>
    <div class="container">
        {{block "yield" .}}{{end}}
    </div>
</body>
</html>`

	// Example 2: Custom login form template with modern styling
	customLogin := `{{block "title" .}}Welcome to {{.ProductName}}{{end}}

{{define "yield"}}
<div style="padding: 40px 30px;">
    <div style="text-align: center; margin-bottom: 30px;">
        <img src="{{.Logo}}" alt="{{.ProductName}}" style="height: 50px; margin-bottom: 20px;">
        <h1 style="margin: 0; color: #333; font-size: 24px; font-weight: 600;">Welcome Back</h1>
        <p style="margin: 10px 0 0 0; color: #666; font-size: 14px;">Sign in to your {{.ProductName}} account</p>
    </div>

    {{if ne .Error ""}}
    <div style="background: #fee; border: 1px solid #fcc; color: #a00; padding: 12px; border-radius: 6px; margin-bottom: 20px; font-size: 14px;">
        <strong>Error:</strong> {{.Error}}
    </div>
    {{end}}

    <form action="{{prefixedPath "/email"}}" method="POST" style="margin: 0;">
        <div style="margin-bottom: 20px;">
            <label for="email" style="display: block; margin-bottom: 8px; color: #333; font-weight: 500; font-size: 14px;">
                Email Address
            </label>
            <input 
                type="email" 
                id="email" 
                name="email" 
                required 
                autofocus
                placeholder="Enter your email address"
                style="width: 100%; padding: 12px 16px; border: 2px solid #e1e5e9; border-radius: 8px; font-size: 16px; box-sizing: border-box; transition: border-color 0.2s;"
                onfocus="this.style.borderColor='#667eea'"
                onblur="this.style.borderColor='#e1e5e9'"
            >
        </div>
        
        <button 
            type="submit"
            style="width: 100%; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; border: none; padding: 14px; border-radius: 8px; font-size: 16px; font-weight: 600; cursor: pointer; transition: transform 0.2s;"
            onmouseover="this.style.transform='translateY(-1px)'"
            onmouseout="this.style.transform='translateY(0)'"
        >
            Send Login Code
        </button>
    </form>

    <div style="text-align: center; margin-top: 20px; padding-top: 20px; border-top: 1px solid #eee;">
        <p style="margin: 0; color: #666; font-size: 12px;">
            Secure login powered by {{.ProductName}}
        </p>
    </div>
</div>
{{end}}`

	// Example 3: Minimal custom template
	minimalLogin := `{{define "yield"}}
<div style="padding: 20px; max-width: 300px; margin: 50px auto; border: 1px solid #ddd; border-radius: 8px;">
    <h2>{{.ProductName}} Login</h2>
    {{if ne .Error ""}}<p style="color: red;">{{.Error}}</p>{{end}}
    <form action="{{prefixedPath "/email"}}" method="POST">
        <input type="email" name="email" placeholder="Email" required style="width: 100%; padding: 10px; margin-bottom: 10px;">
        <button type="submit" style="width: 100%; padding: 10px; background: #007cba; color: white; border: none;">Login</button>
    </form>
</div>
{{end}}`

	fmt.Println("Starting Maildoor examples...")

	// Example server with custom layout and login templates
	handler1 := maildoor.New(
		maildoor.CustomLayoutTemplate(customLayout),
		maildoor.CustomLoginTemplate(customLogin),
		maildoor.ProductName("CustomApp"),
		maildoor.Logo("https://via.placeholder.com/150x50/667eea/white?text=CustomApp"),
		maildoor.Prefix("/auth"),
		maildoor.EmailSender(func(to, html, txt string) error {
			fmt.Printf("Sending email to: %s\n", to)
			return nil
		}),
		maildoor.AfterLogin(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "<h1>Welcome! You are now logged in.</h1>")
		}),
	)

	// Example server with minimal custom template
	handler2 := maildoor.New(
		maildoor.CustomLoginTemplate(minimalLogin),
		maildoor.ProductName("MinimalApp"),
		maildoor.Prefix("/simple"),
		maildoor.EmailSender(func(to, html, txt string) error {
			fmt.Printf("Minimal app - sending email to: %s\n", to)
			return nil
		}),
	)

	// Example server with only custom layout (using default login form)
	handler3 := maildoor.New(
		maildoor.CustomLayoutTemplate(customLayout),
		maildoor.ProductName("HybridApp"),
		maildoor.Prefix("/hybrid"),
		maildoor.EmailSender(func(to, html, txt string) error {
			fmt.Printf("Hybrid app - sending email to: %s\n", to)
			return nil
		}),
	)

	mux := http.NewServeMux()
	mux.Handle("/auth/", handler1)
	mux.Handle("/simple/", handler2)
	mux.Handle("/hybrid/", handler3)

	// Add a home page showing the different examples
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
<html>
<head><title>Maildoor Custom Template Examples</title></head>
<body style="font-family: sans-serif; padding: 40px; max-width: 600px; margin: 0 auto;">
    <h1>Maildoor Custom Template Examples</h1>
    <p>Click on the links below to see different custom template implementations:</p>
    
    <div style="margin: 20px 0;">
        <h3>1. Full Custom Layout + Login Template</h3>
        <p>Modern gradient design with custom styling</p>
        <a href="/auth/login" style="color: #667eea; text-decoration: none; font-weight: bold;">→ View Custom Template</a>
    </div>
    
    <div style="margin: 20px 0;">
        <h3>2. Minimal Custom Login Template</h3>
        <p>Simple, clean design with basic styling</p>
        <a href="/simple/login" style="color: #667eea; text-decoration: none; font-weight: bold;">→ View Minimal Template</a>
    </div>
    
    <div style="margin: 20px 0;">
        <h3>3. Custom Layout + Default Login</h3>
        <p>Custom layout with the default Maildoor login form</p>
        <a href="/hybrid/login" style="color: #667eea; text-decoration: none; font-weight: bold;">→ View Hybrid Template</a>
    </div>
    
    <hr style="margin: 40px 0;">
    
    <h3>How to Test:</h3>
    <ol>
        <li>Click on any of the login pages above</li>
        <li>Enter any email address (e.g., test@example.com)</li>
        <li>Check the console output to see the "email" being sent</li>
        <li>Note the different styling and layouts</li>
    </ol>
    
    <h3>Key Features Demonstrated:</h3>
    <ul>
        <li>Custom HTML layout templates</li>
        <li>Custom login form templates</li>
        <li>CSS styling within templates</li>
        <li>Error message handling</li>
        <li>Form action URL generation with prefixed paths</li>
        <li>Template data binding (Logo, ProductName, Error)</li>
        <li>Mixed custom and default templates</li>
    </ul>
</body>
</html>`)
	})

	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("Visit http://localhost:8080 to see the examples")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
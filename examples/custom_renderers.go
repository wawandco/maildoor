package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wawandco/maildoor"
)

func main() {
	// Example of using custom renderers with maildoor
	handler := maildoor.New(
		maildoor.ProductName("MyApp"),
		maildoor.Logo("https://example.com/logo.png"),
		maildoor.LoginRenderer(customLoginRenderer),
		maildoor.CodeRenderer(customCodeRenderer),
		maildoor.EmailSender(mockEmailSender),
	)

	fmt.Println("Server starting on :8080")
	fmt.Println("Visit http://localhost:8080/login to see custom renderers in action")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// customLoginRenderer renders a custom login page
func customLoginRenderer(data maildoor.Attempt) (string, error) {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Custom Login</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            margin: 0;
            padding: 0;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .container {
            background: white;
            padding: 2rem;
            border-radius: 10px;
            box-shadow: 0 10px 25px rgba(0,0,0,0.1);
            width: 100%%;
            max-width: 400px;
        }
        .logo {
            text-align: center;
            margin-bottom: 2rem;
        }
        .logo img {
            height: 60px;
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 1rem;
        }
        .description {
            color: #666;
            text-align: center;
            margin-bottom: 2rem;
            font-size: 14px;
        }
        .form-group {
            margin-bottom: 1rem;
        }
        label {
            display: block;
            margin-bottom: 0.5rem;
            color: #333;
            font-weight: bold;
        }
        input[type="email"] {
            width: 100%%;
            padding: 12px;
            border: 2px solid #ddd;
            border-radius: 5px;
            font-size: 16px;
            box-sizing: border-box;
        }
        input[type="email"]:focus {
            border-color: #667eea;
            outline: none;
        }
        .error {
            color: #e74c3c;
            font-size: 14px;
            margin-top: 0.5rem;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }
        .submit-btn {
            width: 100%%;
            padding: 12px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            border: none;
            border-radius: 5px;
            font-size: 16px;
            cursor: pointer;
            transition: transform 0.2s;
        }
        .submit-btn:hover {
            transform: translateY(-1px);
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <img src="%s" alt="%s Logo">
        </div>
        <h1>Welcome to %s</h1>
        <p class="description">
            Enter your email address and we'll send you a secure login code.
        </p>
        <form action="/email" method="POST">
            <div class="form-group">
                <label for="email">Email Address</label>
                <input type="email" id="email" name="email" placeholder="your@email.com" required autofocus>
                %s
            </div>
            <button type="submit" class="submit-btn">Send Login Code</button>
        </form>
    </div>
</body>
</html>`, data.ProductName, data.Logo, data.ProductName, data.ProductName, renderError(data.Error))

	return html, nil
}

// customCodeRenderer renders a custom code entry page
func customCodeRenderer(data maildoor.Attempt) (string, error) {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Enter Code</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            margin: 0;
            padding: 0;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .container {
            background: white;
            padding: 2rem;
            border-radius: 10px;
            box-shadow: 0 10px 25px rgba(0,0,0,0.1);
            width: 100%%;
            max-width: 400px;
            text-align: center;
        }
        .logo {
            margin-bottom: 2rem;
        }
        .logo img {
            height: 60px;
        }
        h1 {
            color: #333;
            margin-bottom: 1rem;
        }
        .description {
            color: #666;
            margin-bottom: 2rem;
            font-size: 14px;
            line-height: 1.5;
        }
        .email-highlight {
            color: #667eea;
            font-weight: bold;
        }
        .form-group {
            margin-bottom: 1.5rem;
        }
        .code-input {
            width: 100%%;
            padding: 15px;
            border: 2px solid #ddd;
            border-radius: 5px;
            font-size: 24px;
            text-align: center;
            letter-spacing: 8px;
            font-family: monospace;
            box-sizing: border-box;
        }
        .code-input:focus {
            border-color: #667eea;
            outline: none;
        }
        .error {
            color: #e74c3c;
            font-size: 14px;
            margin-top: 0.5rem;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 0.5rem;
        }
        .submit-btn {
            width: 100%%;
            padding: 12px;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            border: none;
            border-radius: 5px;
            font-size: 16px;
            cursor: pointer;
            transition: transform 0.2s;
            margin-bottom: 1rem;
        }
        .submit-btn:hover {
            transform: translateY(-1px);
        }
        .back-link {
            color: #667eea;
            text-decoration: none;
            font-size: 14px;
        }
        .back-link:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <img src="%s" alt="%s Logo">
        </div>
        <h1>Check Your Email</h1>
        <p class="description">
            We've sent a 6-digit code to<br>
            <span class="email-highlight">%s</span><br><br>
            Enter the code below to continue.
        </p>
        <form action="/code" method="POST">
            <input type="hidden" name="email" value="%s">
            <div class="form-group">
                <input type="text" name="code" class="code-input" maxlength="6" placeholder="000000" required autofocus>
                %s
            </div>
            <button type="submit" class="submit-btn">Verify Code</button>
        </form>
        <a href="/login" class="back-link">← Back to login</a>
    </div>
</body>
</html>`, data.ProductName, data.Logo, data.ProductName, data.Email, data.Email, renderError(data.Error))

	return html, nil
}

// renderError creates HTML for error messages
func renderError(errorMsg string) string {
	if errorMsg == "" {
		return ""
	}
	return fmt.Sprintf(`<div class="error">⚠️ %s</div>`, errorMsg)
}

// mockEmailSender simulates sending emails (for demo purposes)
func mockEmailSender(to, html, txt string) error {
	fmt.Printf("📧 Mock email sent to: %s\n", to)
	fmt.Printf("📧 Email content: %s\n", txt)
	return nil
}
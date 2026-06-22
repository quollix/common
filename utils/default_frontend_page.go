package utils

import "fmt"

func RenderDefaultPage(htmlContent string) string {
	return fmt.Sprintf(defaultPageTemplate, htmlContent)
}

const defaultPageTemplate = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Server</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
body{margin:0;font-family:system-ui,-apple-system,Segoe UI,Roboto,Ubuntu,Cantarell,Noto Sans,sans-serif;background:#0b1220;color:#e6e9ef;display:flex;align-items:center;justify-content:center;min-height:100vh}
.card{max-width:720px;padding:32px;border-radius:16px;background:#121a2b;box-shadow:0 10px 30px rgba(0,0,0,.35);text-align:center}
h1{margin:0 0 8px;font-size:28px}
p{margin:8px 0 0;line-height:1.6;color:#c0c6d4}
a{color:#8dd0ff;text-decoration:none}
a:hover{text-decoration:underline}
</style>
</head>
<body>
  <div class="card">
	%s
  </div>
</body>
</html>`

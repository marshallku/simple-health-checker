package templates

import "html/template"

var IndexTemplate = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Statusy - Current Status</title>
    <meta http-equiv="refresh" content="30">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .status-card {
            border: 1px solid #ddd;
            padding: 10px;
            margin: 10px 0;
            border-radius: 4px;
        }
        .UP { background-color: #d4edda; }
        .DOWN { background-color: #f8d7da; }
        nav { margin-bottom: 20px; }
        nav a { margin-right: 10px; }
    </style>
</head>
<body>
    <nav>
        <a href="/">Status</a>
        <a href="/history">History</a>
    </nav>
    <h1>Current Status</h1>
    {{range .}}
    <div class="status-card {{.Status}}">
        <h3>{{.URL}}</h3>
        <p>Status: {{.Status}}</p>
        <p>Status Code: {{.StatusCode}}</p>
        <p>Response Time: {{.TimeTaken}}</p>
        <p>Last Checked: {{.LastChecked.Format "2006-01-02 15:04:05"}}</p>
    </div>
    {{end}}
</body>
</html>
`))

var HistoryTemplate = template.Must(template.New("history").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Statusy - History</title>
    <meta http-equiv="refresh" content="30">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .history-item {
            border: 1px solid #ddd;
            padding: 10px;
            margin: 10px 0;
            border-radius: 4px;
        }
        nav { margin-bottom: 20px; }
        nav a { margin-right: 10px; }
    </style>
</head>
<body>
    <nav>
        <a href="/">Status</a>
        <a href="/history">History</a>
    </nav>
    <h1>History (Last 10 Events)</h1>
    {{range .}}
    <div class="history-item">
        <h3>{{.URL}}</h3>
        <p>Status: {{.Status}}</p>
        <p>Message: {{.Message}}</p>
        <p>Time: {{.Timestamp.Format "2006-01-02 15:04:05"}}</p>
    </div>
    {{end}}
</body>
</html>
`))

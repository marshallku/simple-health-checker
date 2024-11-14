package templates

import "html/template"

var IndexTemplate = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Statusy - Current Status</title>
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
    <div id="status-container"></div>

    <script>
        const statusContainer = document.getElementById('status-container');

        function updateStatus(results) {
            statusContainer.innerHTML = Object.values(results)
                .map(result => ` + "`" + `
                    <div class="status-card ${result.Status}">
                        <h3>${result.url}</h3>
                        <p>Status: ${result.status ? 'UP' : 'DOWN'}</p>
                        <p>Status Code: ${result.statusCode}</p>
                        <p>Response Time: ${result.timeTaken}</p>
                        <p>Last Checked: ${new Date(result.lastChecked).toLocaleString()}</p>
                    </div>
                ` + "`" + `).join('');
        }

        // Connect to WebSocket
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const ws = new WebSocket(protocol + '//' + window.location.host + '/ws');

        ws.onmessage = function(event) {
            const data = JSON.parse(event.data);
            if (data.type !== 'results') return;
            updateStatus(data.data);
        };

        ws.onclose = function() {
            setTimeout(() => {
                window.location.reload();
            }, 1000);
        };
    </script>
</body>
</html>
`))

var HistoryTemplate = template.Must(template.New("history").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Statusy - History</title>
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
    <div id="history-container"></div>

    <script>
        const historyContainer = document.getElementById('history-container');

        function updateHistory(history) {
            historyContainer.innerHTML = history
                .map(item => ` + "`" + `
                    <div class="history-item">
                        <h3>${item.url}</h3>
                        <p>Status: ${item.status}</p>
                        <p>Time: ${new Date(item.timestamp).toLocaleString()}</p>
                    </div>
                ` + "`" + `).join('');
        }

        // Connect to WebSocket
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const ws = new WebSocket(protocol + '//' + window.location.host + '/ws');

        ws.onmessage = function(event) {
            const data = JSON.parse(event.data);
            if (data.type !== 'history') return;
            updateHistory(data.data);
        };

        ws.onclose = function() {
            setTimeout(() => {
                window.location.reload();
            }, 1000);
        };
    </script>
</body>
</html>
`))

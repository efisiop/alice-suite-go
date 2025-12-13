const http = require('http');
const fs = require('fs');

function createHelpRequest() {
    const postData = JSON.stringify({
        user_id: 'test_user',
        book_id: 'alice-in-wonderland',
        content: 'Test help request',
        context: 'Test context',
        section_id: null
    });

    const options = {
        hostname: 'localhost',
        port: 8080,
        path: '/api/help/request',
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Content-Length': Buffer.byteLength(postData)
        }
    };

    const req = http.request(options, (res) => {
        console.log(`STATUS: ${res.statusCode}`);
        console.log(`HEADERS: ${JSON.stringify(res.headers)}`);
        
        let body = '';
        res.on('data', (chunk) => {
            body += chunk;
        });
        res.on('end', () => {
            console.log(`BODY: ${body}`);
        });
    });

    req.on('error', (e) => {
        console.error(`problem with request: ${e.message}`);
    });

    req.write(postData);
    req.end();
}

// Start server and test
go run ./cmd/reader 2>&1 > server.log &
echo "Server started with PID $!"
sleep 2
createHelpRequest
sleep 1
kill $! 2>&1

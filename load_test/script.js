import { check } from "k6";
import http from "k6/http";

export const options = {
    stages: [
        {
            duration: '1m',
            target: 1000
        },
        {
            duration: '30s',
            target: 0
        }
    ],
    thresholds: {
        http_req_failed: ['rate<0.001'], // the error rate must be lower than 0.1%
        http_req_duration: ['p(90)<600'], // 90% of requests must complete below 600ms
        http_req_receiving: ['max<150'], // slowest request below 150ms
        iteration_duration: ['p(95)<20000'], // 95% of requests must complete below 20000ms
    },
};

export default function() {
    let body = {
        "original": "https://www.google.com/",
        "expiry": 60
    }
    let res = http.post("http://localhost:8080/shorten", JSON.stringify(body));
    check(res, {
        "is status 200": (r) => r.status === 200
    });
};

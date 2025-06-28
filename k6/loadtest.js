import http from 'k6/http';
import { sleep, check } from 'k6';

export let options = {
    vus: 100,
    duration: '30s',
};

export default function () {
    // Get all users
    let res = http.get('http://api-bad:8080/users');
    check(res, {
        'status is 200': (r) => r.status === 200,
        'body is not empty': (r) => r.body.length > 0,
    });

    // Get user by ID with random ID
    let userId = Math.floor(Math.random() * 90000) + 1; // Assuming IDs are 1-90000
    res = http.get(`http://api-bad:8080/users/${userId}`);
    check(res, {
        'status is 200': (r) => r.status === 200,
        'body is not empty': (r) => r.body.length > 0,
    });

}

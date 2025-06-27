import http from 'k6/http';
import { sleep, check } from 'k6';

export let options = {
    vus: 100,
    duration: '30s',
};

export default function () {
    let res = http.get('http://api-bad:8080/users');
    check(res, {
        'status is 200': (r) => r.status === 200,
        'body is not empty': (r) => r.body.length > 0,
    });
    sleep(1);
}

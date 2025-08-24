import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';

export const options = {
  vus: 10,
  duration: '30s',
};

// bikin counter per status
const count202 = new Counter('transfer_202');
const count400 = new Counter('transfer_400');
const count500 = new Counter('transfer_500');

export default function () {
  const url = 'http://doitpay-api:80/transfer';
  const payload = JSON.stringify({
    amount: 10000,
    bank_code: 'BCA',
    beneficiary_account: '12345',
    beneficiary_name: 'Test Name',
    note: 'coba transfer',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'accept': 'application/json',
    },
  };

  let res = http.post(url, payload, params);

  // increment counter sesuai status
  if (res.status === 202) {
    count202.add(1);
  } else if (res.status === 400) {
    count400.add(1);
  } else if (res.status === 500) {
    count500.add(1);
  }

  check(res, {
    'is status 202/400/500': (r) => [202, 400, 500].includes(r.status),
  });

  sleep(1);
}

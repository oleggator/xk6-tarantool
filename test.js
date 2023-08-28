import tarantool from 'k6/x/tarantool';

const tarantool_addrs = __ENV.TARANTOOL_ADDRS || 'localhost:3301';
const tarantool_user = __ENV.TARANTOOL_USER || '';
const tarantool_password = __ENV.TARANTOOL_PASSWORD || '';

const conn = new tarantool.Client({
  addrs: tarantool_addrs.split(','),
  user: tarantool_user,
  password: tarantool_password,
});

export const setup = () => {
  console.log('setup');
};

export default async () => {
  const resp = await conn.call('test_func', ['some value']);
  console.log(resp)
};

export const teardown = () => {
  console.log('teardown');
};

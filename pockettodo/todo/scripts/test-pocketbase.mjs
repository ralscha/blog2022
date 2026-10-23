import assert from 'node:assert/strict';
import { spawn, spawnSync } from 'node:child_process';
import { once } from 'node:events';
import { copyFile, mkdir, mkdtemp, readFile, rm } from 'node:fs/promises';
import { createServer } from 'node:net';
import { tmpdir } from 'node:os';
import { basename, dirname, join, resolve } from 'node:path';
import { setTimeout as delay } from 'node:timers/promises';
import { fileURLToPath } from 'node:url';
import PocketBase from 'pocketbase';

const project = fileURLToPath(new URL('../', import.meta.url));
const backend = resolve(project, '..');
const binary =
  process.env.POCKETBASE_BIN ||
  join(backend, process.platform === 'win32' ? 'pocketbase.exe' : 'pocketbase');
const version = (
  await readFile(join(backend, 'pocketbase.version'), 'utf8')
).trim();
const reported = spawnSync(binary, ['--version'], {
  encoding: 'utf8',
  windowsHide: true
});
if (reported.error) throw reported.error;
assert.equal(reported.status, 0, reported.stderr);
assert.equal(
  reported.stdout.trim().split(' ').at(-1),
  version,
  'Unexpected PocketBase version'
);

// Reserve an available port briefly; no application database or server is used.
const reservation = createServer();
reservation.listen(0, '127.0.0.1');
await once(reservation, 'listening');
const port = reservation.address().port;
await new Promise((resolve, reject) =>
  reservation.close(error => (error ? reject(error) : resolve()))
);
const address = `http://127.0.0.1:${port}`;
const temporary = await mkdtemp(join(tmpdir(), 'pockettodo-test-'));
const previousBinary = process.env.POCKETBASE_PREVIOUS_BIN;
let activeBinary = previousBinary || binary;
let migrations = join(backend, 'pb_migrations');
let server;
let output = '';

async function start() {
  server = spawn(
    activeBinary,
    [
      'serve',
      `--http=127.0.0.1:${port}`,
      `--dir=${join(temporary, 'pb_data')}`,
      `--migrationsDir=${migrations}`,
      '--automigrate=false'
    ],
    { windowsHide: true, stdio: ['ignore', 'pipe', 'pipe'] }
  );
  server.stdout.on('data', chunk => {
    output += chunk;
  });
  server.stderr.on('data', chunk => {
    output += chunk;
  });
  let spawnError;
  server.on('error', error => {
    spawnError = error;
  });
  for (let attempt = 0; attempt < 100; attempt++) {
    if (spawnError) throw spawnError;
    assert.equal(server.exitCode, null, output);
    try {
      const response = await fetch(`${address}/api/health`, {
        signal: AbortSignal.timeout(1000)
      });
      if (response.ok) return;
    } catch {
      /* Wait for startup and migrations. */
    }
    await delay(100);
  }
  throw new Error(`PocketBase startup timed out:\n${output}`);
}

async function stop() {
  if (server && server.exitCode === null && server.signalCode === null) {
    const exited = once(server, 'exit');
    server.kill();
    await exited;
  }
}

try {
  if (previousBinary) {
    // Reproduce the original schema before applying the later ownership rule.
    migrations = join(temporary, 'pb_migrations');
    await mkdir(migrations);
    for (const name of [
      '1752857221_created_todos.js',
      '1752857440_updated_todos.js'
    ]) {
      await copyFile(
        join(backend, 'pb_migrations', name),
        join(migrations, name)
      );
    }
  }
  // Precreate a disposable superuser so startup doesn't open an installer browser.
  const setup = spawnSync(
    activeBinary,
    [
      'superuser',
      'create',
      'admin@example.com',
      'temporary-test-password',
      `--dir=${join(temporary, 'pb_data')}`,
      `--migrationsDir=${migrations}`
    ],
    { encoding: 'utf8', windowsHide: true }
  );
  if (setup.error) throw setup.error;
  assert.equal(setup.status, 0, setup.stderr);
  await start();

  const owner = new PocketBase(address);
  const other = new PocketBase(address);
  const anonymous = new PocketBase(address);
  const password = 'temporary-test-password';
  async function register(client, email) {
    const user = await client
      .collection('users')
      .create({ email, password, passwordConfirm: password });
    const auth = await client
      .collection('users')
      .authWithPassword(email, password);
    assert.equal(auth.record.id, user.id);
    assert.equal(
      (await client.collection('users').authRefresh()).record.id,
      user.id
    );
    return user;
  }
  const alice = await register(owner, 'alice@example.com');
  const bob = await register(other, 'bob@example.com');
  const todo = await owner.collection('todos').create({
    title: 'Upgrade PocketBase',
    description: 'Check the current API',
    due_date: '2026-10-23 12:00:00.000Z',
    user: alice.id
  });
  if (previousBinary) {
    await stop();
    activeBinary = binary;
    migrations = join(backend, 'pb_migrations');
    await start();
    // Existing users, tokens, and todos must survive the server/schema upgrade.
    assert.equal(
      (await owner.collection('users').authRefresh()).record.id,
      alice.id
    );
    assert.equal(
      (await other.collection('users').authRefresh()).record.id,
      bob.id
    );
  }
  assert.equal(
    (await owner.collection('todos').getOne(todo.id)).title,
    todo.title
  );
  assert.equal(
    (await owner.collection('todos').getFullList({ sort: '-created' })).length,
    1
  );
  assert.equal((await other.collection('todos').getFullList()).length, 0);
  assert.equal((await anonymous.collection('todos').getFullList()).length, 0);
  await assert.rejects(other.collection('todos').getOne(todo.id), {
    status: 404
  });
  await assert.rejects(
    other.collection('todos').update(todo.id, { title: 'Forbidden' }),
    { status: 404 }
  );
  await assert.rejects(other.collection('todos').delete(todo.id), {
    status: 404
  });
  await assert.rejects(
    anonymous
      .collection('todos')
      .create({ title: 'Forbidden', user: alice.id }),
    { status: 400 }
  );
  await assert.rejects(
    other.collection('todos').create({ title: 'Forged owner', user: alice.id }),
    { status: 400 }
  );
  await assert.rejects(
    owner.collection('todos').update(todo.id, { user: bob.id }),
    { status: 404 }
  );
  await owner.collection('todos').update(todo.id, { completed: true });
  assert.equal(
    (
      await owner
        .collection('todos')
        .getFullList({ filter: 'completed = false' })
    ).length,
    0
  );
  assert.equal(
    (await owner.collection('todos').getOne(todo.id)).user,
    alice.id
  );

  await stop();
  await start();
  assert.equal(
    (await owner.collection('todos').getOne(todo.id)).completed,
    true
  );
  await owner.collection('todos').delete(todo.id);
  assert.equal((await owner.collection('todos').getFullList()).length, 0);
  console.log(
    `PocketBase ${version}: ${previousBinary ? 'upgrade, ' : ''}migrations, authentication, CRUD, access rules, and restart passed.`
  );
} catch (error) {
  console.error(output);
  throw error;
} finally {
  await stop();
  assert.equal(dirname(temporary), resolve(tmpdir()));
  assert.ok(basename(temporary).startsWith('pockettodo-test-'));
  await rm(temporary, { recursive: true, force: true });
}

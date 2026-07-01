import { expect, test } from 'bun:test';
import type { HeartbeatPayload } from '../heartbeat';
import { HeartbeatQueue } from '../queue';

function makePayload(): HeartbeatPayload {
  return {
    file: '/test/file.ts',
    project: 'test-project',
    language: 'typescript',
    editor: 'vscode',
    time: Date.now(),
  };
}

test('enqueue adds a heartbeat to the queue', () => {
  const queue = new HeartbeatQueue();
  queue.enqueue(makePayload());
  expect(queue.size()).toBe(1);
});

test('enqueue drops the oldest item when queue exceeds max size', () => {
  const queue = new HeartbeatQueue();
  for (let i = 0; i < 105; i++) {
    queue.enqueue(makePayload());
  }
  expect(queue.size()).toBe(100);
});

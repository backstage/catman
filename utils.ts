/*
 * Copyright 2021 Spotify AB
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

const globalTimers = new Map<string, number[]>();

function now() {
  const [s, ns] = process.hrtime();
  return s + ns / 1000_000_000;
}

export function createTimer(name: string) {
  const start = now();
  return () => {
    const duration = now() - start;
    let durations = globalTimers.get(name);
    if (!durations) {
      durations = [];
      globalTimers.set(name, durations);
    }
    durations!.push(duration);
  };
}

setInterval(() => {
  console.log('### TIMER SUMMARY ###');
  globalTimers.forEach((durations, name) => {
    if (durations.length) {
      const avg = durations.reduce((sum, x) => sum + x, 0) / durations.length;
      console.log(
        `Timer ${name}(${(durations.length / 5).toFixed(0)}): ${(
          avg * 1000
        ).toFixed(1)} ms`,
      );
      durations.length = 0;
    }
  });
}, 5000);

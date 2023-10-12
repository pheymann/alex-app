
export function logError({ awsFetch, error, logEntriesRef }) {
  const message = {
    message: error.message,
    stack: error.stack,
    name: error.name,
    cause: error.cause,
  };

  pushLogMessage(logEntriesRef, { level: 'error', message: JSON.stringify(message)});

  const logEntries = logEntriesRef.current;

  awsFetch.call(`/api/app/logs`, {
    method: 'POST',
    body: JSON.stringify({
      logEntries,
    }),
  });
}

export function pushLogMessage(logEntriesRef, { level, message }) {
  logEntriesRef.current.push({
    level,
    timestamp: new Date().toISOString(),
    message,
  });
}

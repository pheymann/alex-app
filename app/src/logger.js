
export function logError({ token, error, logEntriesRef }) {
  pushLogMessage(logEntriesRef, { level: 'error', message: error.message });

  const logEntries = logEntriesRef.current;

  fetch(`/api/app/logs`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
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

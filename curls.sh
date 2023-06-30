curl 'https://app.resemble.ai/api/v2/projects/bc7a8b34/clips' \
  -H 'Authorization: Token token=jHPsa2ctfTx0xxwYEU8zjAtt' \
  -H 'Content-Type: application/json' \
  --data '{"title": "CLI Test", "body": "This audio was synthesized", "voice_uuid": "7c8e47ca", "is_public": false, "is_archived": false, "raw": true}'

curl 'https://app.resemble.ai/api/v2/projects/bc7a8b34/clips' \
  -H 'Authorization: Token token=jHPsa2ctfTx0xxwYEU8zjAtt' \
  -H 'Content-Type: application/json' \
  --data '{"callback_uri": "https://test.me", "title": "CLI Test", "body": "This audio was synthesized", "voice_uuid": "7c8e47ca", "is_public": false, "is_archived": false, "raw": true}'


curl 'https://app.resemble.ai/api/v2/voices?page=1&page_size=10' \
  -H 'Authorization: Token token=jHPsa2ctfTx0xxwYEU8zjAtt'


curl 'https://app.resemble.ai/api/v2/projects/bc7a8b34/clips/1b836d55' \
  -H 'Authorization: Token token=jHPsa2ctfTx0xxwYEU8zjAtt'


curl -vL  'https://app.resemble.ai/rails/active_storage/blobs/redirect/eyJfcmFpbHMiOnsibWVzc2FnZSI6IkJBaHBCTUFNQ3cwPSIsImV4cCI6bnVsbCwicHVyIjoiYmxvYl9pZCJ9fQ==--b9d5ed35f6bf3f3deb0f5e44e9ad77b19ec8203b/CLI+Test-238771b5.wav' \
  -H 'Authorization: Token token=jHPsa2ctfTx0xxwYEU8zjAtt' \
  --output test.wav


curl 'https://app.resemble.ai/api/v2/projects/bc7a8b34/clips?page=1&page_size=10' \
  -H 'Authorization: Token token=jHPsa2ctfTx0xxwYEU8zjAtt'

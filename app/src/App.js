import './App.css';
import React, { useState } from 'react';

function App() {
  const [artistNames, setArtistName] = useState('');
  const [artPieceName, setArtPieceName] = useState('');

  const [messages, setMessages] = useState([]);
  const [audioClip, setAudioClip] = useState('');

  const handleClick = () => {
    fetch(`/api/art`, {
      method: 'POST',
      body: JSON.stringify({
        artist_name: artistNames,
        art_piece_name: artPieceName,
      }),
    })
      .then(response => response.json())
      .then(data => {
        console.log(data);
        setMessages(data.conversation_start.messages);
        setAudioClip(data.conversation_start_clip_uuid);
      })
      .catch(error => {
        console.log(error);
      });
  };

  return (
    <div>
      <input
        type="text"
        value={artistNames}
        onChange={(e) => setArtistName(e.target.value)}
      />
      <input
        type="text"
        value={artPieceName}
        onChange={(e) => setArtPieceName(e.target.value)}
      />
      <button onClick={handleClick}>Send Request</button>

      { messages &&
        <div>
          {
            messages.map((message, index) => {
              return <p key={index}>{message.text}</p>;
            })
          }

          {audioClip && <audio src="/api/assets/{audioClip}" controls />}
        </div>
      }
    </div>
  );
}

export default App;

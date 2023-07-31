import './App.css';
import React, { useState } from 'react';

function App() {
  const [artistNames, setArtistName] = useState('');
  const [artPieceName, setArtPieceName] = useState('');

  const [messages, setMessages] = useState([]);

  const handleStartConversation = () => {
    fetch(`/api/conversation/create/art`, {
      method: 'POST',
      body: JSON.stringify({
        artistName: artistNames,
        artPiece: artPieceName,
        userUuid: "1",
      }),
    })
      .then(response => response.json())
      .then(data => {
        console.log(data);
        setMessages(data.messages);
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
      <button onClick={handleStartConversation}>Start Conversation</button>

      { messages &&
        <div>
          {
            messages.map((message, index) => {
              return <div key="{message.speechClipUuid}">
                  <audio src="/api/assets/{message.speechClipUuid}" controls />
                  <p key={index}>{message.text}</p>
                </div>
            })
        }
        </div>
      }
    </div>
  );
}

export default App;

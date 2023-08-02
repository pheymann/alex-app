import './App.css';
import React, { useState } from 'react';

function App() {
  const [artistNames, setArtistName] = useState('');
  const [artPieceName, setArtPieceName] = useState('');

  const [conversation, setConversation] = useState(null);
  const [prompt, setPrompt] = useState('');

  const getRandomInt = () => {
    return Math.floor(Math.random() * 101);
  }

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
        setConversation(data);
      })
      .catch(error => {
        console.log(error);
      });
  };

  const handlePrompt = () => {
    fetch(`/api/conversation/continue`, {
      method: 'POST',
      body: JSON.stringify({
        conversationUuid: conversation.id,
        userUuid: "1",
        prompt: prompt,
      }),
    })
      .then(response => response.json())
      .then(data => {
        const newConversation = {
          ...conversation,
          messages: [...conversation.messages,
            {
              text: prompt,
            },
            {
              text: data.text,
              speechClipUuid: data.speechClipUuid,
            }],
        };

        setConversation(newConversation);
        setPrompt('');
      })
      .catch(error => {
        console.log(error);
      });
  }

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
      { !conversation &&
        <button onClick={handleStartConversation}>Start Conversation</button>
      }

      { conversation &&
        <div>
          {

            conversation.messages.map((message, index) => {
              const key = message.speechClipUuid ? message.speechClipUuid : getRandomInt();
              return <div key={key}>
                  {message.speechClipUuid && <audio src={'/api/assets/' + message.speechClipUuid} controls /> }
                  <p>{message.text}</p>
                </div>
            })
          }
          <div>
            <input
              type="text"
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
            />
             <button onClick={handlePrompt}>Ask</button>
          </div>
        </div>
      }
    </div>
  );
}

export default App;

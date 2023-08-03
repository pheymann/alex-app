import React, { useEffect, useState } from 'react';

export default function ArtConversation(conversationId) {
  const [artistNames, setArtistName] = useState('');
  const [artPieceName, setArtPieceName] = useState('');

  const [conversation, setConversation] = useState(null);
  const [prompt, setPrompt] = useState('');

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

  useEffect(() => {
    if (!conversationId) {
      return;
    }

    fetch(`/api/conversation/${conversationId}`, {
      method: 'POST',
      body: JSON.stringify({
        userUuid: "1",
      }),
    })
      .then(response => response.json())
      .then(data => {
        setConversation(data);
        setArtistName(data.artistName);
        setArtPieceName(data.artPiece);
      })
      .catch(error => {
        console.log(error);
      });
  }, [conversationId]);

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
              const key = `${message.speechClipUuid}_${index}`;

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

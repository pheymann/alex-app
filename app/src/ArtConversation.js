import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

export default function ArtConversation() {
  const pathParams = useParams();

  const [artistNames, setArtistName] = useState('');
  const [artPieceName, setArtPieceName] = useState('');

  const [conversation, setConversation] = useState(null);
  const [prompt, setPrompt] = useState('');

  const handleStartConversation = () => {
    if (artistNames === '' || artPieceName === '') {
      // TODO: show error message
      console.error('missing artist name or art piece name');
      return;
    }

    fetch(`/api/conversation/create/art`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'User-UUID': '1',
      },
      body: JSON.stringify({
        artistName: artistNames,
        artPiece: artPieceName,
      }),
    })
      .then(response => response.json())
      .then(data => {
        setConversation(data);
      })
      .catch(error => {
        console.log(error);
      });
  };

  const handlePrompt = () => {
    if (prompt === '') {
      // TODO: show error message
      console.error('missing prompt');
      return;
    }

    fetch(`/api/conversation/continue`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'User-UUID': '1',
      },
      body: JSON.stringify({
        conversationUuid: conversation.id,
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

  const conversationId = pathParams.id;

  useEffect(() => {
    if (!conversationId || conversationId === 'new') {
      return;
    }

    fetch(`/api/conversation/${conversationId}`, {
      method: 'GET',
      headers: {
        'User-UUID': '1',
      },
    })
      .then(response => response.json())
      .then(data => {
        setConversation(data);
        setArtistName(data.metadata.artistName);
        setArtPieceName(data.metadata.artPiece);
      })
      .catch(error => {
        console.log(error);
      });
  }, [conversationId]);

  return (
    <div>
      Tell me something about
      <input
        type="text"
        value={artPieceName}
        onChange={(e) => setArtPieceName(e.target.value)}
      />
      by
      <input
        type="text"
        value={artistNames}
        onChange={(e) => setArtistName(e.target.value)}
      />

      { !conversation &&
        <button className="btn btn-primary" onClick={handleStartConversation}>Start Conversation</button>
      }

      { conversation &&
        <div>
          {

            conversation.messages.map((message, index) => {
              const key = `${message.speechClipUuid}_${index}`;

              if (message.role === 'user') {
                return <UserMessage key={key} message={message} />
              } else {
                return <AssistantMessage key={key} index={index} message={message} />
              }
            })
          }
          <div>
            <input
              type="text"
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
            />
             <button className="btn btn-primary" onClick={handlePrompt}>Ask</button>
          </div>
        </div>
      }
    </div>
  );
}

function UserMessage({ message }) {
  return (
    <div className='card card-body'>
      <p>
        {message.text}
      </p>
    </div>
  );
}

function AssistantMessage({ index, message }) {
  const answerInTextId = `answerInText_${index}`;

  return (
    <div>
      <audio src={'/api/assets/speechclip/' + message.speechClipUuid} controls />


      <button className="btn btn-primary" type="button" data-bs-toggle="collapse" data-bs-target={`#${answerInTextId}`} aria-expanded="false">
        Show Text
      </button>
      <div className="collapse" id={answerInTextId}>
        <div className="card card-body">
          {message.text}
        </div>
      </div>
    </div>
  );
}

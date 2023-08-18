import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

export default function ArtConversation({ awsContext }) {
  const pathParams = useParams();

  const [artistName, setArtistName] = useState('');
  const [artPieceName, setArtPieceName] = useState('');
  const [conversation, setConversation] = useState(null);

  const [loading, setLoading] = useState(true);

  const conversationId = pathParams.id;

  useEffect(() => {
    if (!conversationId || conversationId === 'new') {
      setLoading(false);
      return;
    }

    fetch(`/api/conversation/${conversationId}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${awsContext.token}`,
      },
    })
      .then(response => response.json())
      .then(data => {
        setConversation(data);
        setArtistName(data.metadata.artistName);
        setArtPieceName(data.metadata.artPiece);
        setLoading(false);
      })
      .catch(error => {
        console.log(error);
      });
  }, [conversationId, awsContext.token, awsContext.userUUID]);

  if (loading) {
    return(
      <div className="spinner-border" role="status">
        <span className="visually-hidden">Loading...</span>
      </div>
    );
  }

  if (!conversation) {
    return <NewConversation
      artPieceName={artPieceName}
      setArtPieceName={setArtPieceName}
      artistName={artistName}
      setArtistName={setArtistName}
      setConversation={setConversation}
      awsContext={awsContext} />;
  } else {
    return <ContinueConversation
      artPieceName={artPieceName}
      artistName={artistName}
      conversation={conversation}
      setConversation={setConversation}
      awsContext={awsContext} />;
  }
}

function NewConversation({artPieceName, setArtPieceName, artistName, setArtistName, setConversation, awsContext}) {
  const handleStartConversation = () => {
    if (artistName === '' || artPieceName === '') {
      // TODO: show error message
      console.error('missing artist name or art piece name');
      return;
    }

    fetch(`/api/conversation/create/art`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${awsContext.token}`,
      },
      body: JSON.stringify({
        artistName: artistName,
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

  return (
    <div>
      Tell me something about
      <input
        type='text'
        value={artPieceName}
        placeholder='Mona Lisa'
        onChange={(e) => setArtPieceName(e.target.value)}
      />
      by
      <input
        type='text'
        value={artistName}
        placeholder='Leonardo da Vinci'
        onChange={(e) => setArtistName(e.target.value)}
      />

      <button className='btn btn-primary' onClick={handleStartConversation}>Start</button>
    </div>
  );
}

function ContinueConversation({artPieceName, artistName, conversation, setConversation, awsContext}) {
  const [prompt, setPrompt] = useState('');
  const [loading, setLoading] = useState(false);

  const handlePrompt = () => {
    if (prompt === '') {
      // TODO: show error message
      console.error('missing prompt');
      return;
    }

    setLoading(true);

    fetch(`/api/conversation/continue`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${awsContext.token}`,
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
        setLoading(false);
      })
      .catch(error => {
        console.log(error);
      });
  }

  return (
    <div>
      Tell me something about {artPieceName} by {artistName}

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

        { !loading &&
          <div>
            <input
              type='text'
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
            />
            <button className='btn btn-primary' onClick={handlePrompt}>Ask</button>
          </div>
        }

        { loading &&
           <div className="spinner-border" role="status">
             <span className="visually-hidden">Loading...</span>
           </div>
        }
      </div>
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
      <audio src={message.speechClipURL} controls />

      <button className='btn btn-primary' type='button' data-bs-toggle='collapse' data-bs-target={`#${answerInTextId}`} aria-expanded='false'>
        Show Text
      </button>
      <div className='collapse' id={answerInTextId}>
        <div className='card card-body'>
          <p>
            {message.text}
          </p>
        </div>
      </div>
    </div>
  );
}

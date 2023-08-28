import React, { useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import './ArtConversation.css'

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
      <div className="container">
        <div className="spinner-border" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
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
  const [loading, setLoading] = useState(false);

  const handleStartConversation = () => {
    setLoading(true);

    if (artistName === '' || artPieceName === '') {
      // TODO: show error message
      setLoading(false);
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
      .finally(() => {
        setLoading(false);
      })
      .catch(error => {
        console.log(error);
      });
  };

  if (loading) {
    return(
      <div className="container">
        <div className="spinner-border" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    );
  }

  return (
    <div className='container'>
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
  const textareaRef = useRef(null);
  const initialMessage = {
    text: `Tell me something about ${artPieceName} by ${artistName}`,
  };

  const resizeTextArea = () => {
    textareaRef.current.style.height = "auto";
    textareaRef.current.style.height = textareaRef.current.scrollHeight + "px";
  };

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
              role: 'user',
              text: prompt,
            },
            {
              role: data.role,
              text: data.text,
              speechClipUuid: data.speechClipUuid,
              speechClipUrl: data.speechClipUrl,
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

  useEffect(resizeTextArea, [prompt]);

  return (
    <div className='container'>
      <UserMessage key={"init"} message={initialMessage} />

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
          <div className='user-prompt'>
            <textarea
              ref={textareaRef}
              value={prompt}
              rows={1}
              placeholder="Something's on your mind?"
              onChange={(e) => setPrompt(e.target.value)}
            />
            <button onClick={handlePrompt}>
              <span className='arrow'></span>
            </button>
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
    <div className='user-message-bubble'>
      <p>
        {message.text}
      </p>
    </div>
  );
}

function AssistantMessage({ index, message }) {
  const answerInTextId = `answerInText_${index}`;
  const [isCollapsed, setIsCollapsed] = useState(true);

  return (
    <div className='assistant-response-bubble'>
      <audio src={message.speechClipUrl} controls />

      <button className='btn btn-primary'
              type='button'
              data-bs-toggle='collapse'
              data-bs-target={`#${answerInTextId}`}
              aria-expanded='false'
              onClick={_ => setIsCollapsed(!isCollapsed)}>
        {isCollapsed ? 'Show Text' : 'Hide'}
      </button>

      <p className='collapse' id={answerInTextId}>
        {message.text}
      </p>
    </div>
  );
}

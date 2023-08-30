import React, { useEffect, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import './ArtConversation.css'
import NewConversationButton from './NewConversationButton';
import BasicPage from './BasicPage';

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
      setConversation(null);
      setArtPieceName('');
      setArtistName('');
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
      <BasicPage awsContext={awsContext}>
        <div className="container container-limited-width">
          <div className="spinner-border" role="status">
            <span className="visually-hidden">Loading...</span>
          </div>
        </div>
      </BasicPage>
    );
  }

  if (!conversation) {
    return (
      <BasicPage awsContext={awsContext}>
        <NewConversation
          artPieceName={artPieceName}
          setArtPieceName={setArtPieceName}
          artistName={artistName}
          setArtistName={setArtistName}
          awsContext={awsContext} />
      </BasicPage>
    );
  } else {
    return (
      <BasicPage awsContext={awsContext}>
        <ContinueConversation
          artPieceName={artPieceName}
          artistName={artistName}
          conversation={conversation}
          setConversation={setConversation}
          awsContext={awsContext} />
      </BasicPage>
    );
  }
}

function NewConversation({artPieceName, setArtPieceName, artistName, setArtistName, awsContext}) {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

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
        navigate(`/conversation/${data.id}`);
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
      <div className='container container-limited-width'>
        <div className="spinner-border" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    );
  }

  return (
    <div className='container container-limited-width'>
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

    // add user prompt
    const newConversation = {
      ...conversation,
      messages: [...conversation.messages,
        {
          role: 'user',
          text: prompt,
        }],
    };
    setConversation(newConversation);

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
        const responseConversation = {
          ...newConversation,
          messages: [...newConversation.messages,
            {
              role: data.role,
              text: data.text,
              speechClipUuid: data.speechClipUuid,
              speechClipUrl: data.speechClipUrl,
            }],
        };

        setConversation(responseConversation);
        setPrompt('');
        setLoading(false);
      })
      .catch(error => {
        console.log(error);
      });
  }

  useEffect(resizeTextArea, [prompt]);

  return (
    <div className='container container-limited-width'>
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
          <div className='row'>
            <div className='col user-prompt'>
              <div className='row'>
                <div className='col col-10'>
                  <textarea
                    className='user-prompt-textarea'
                    ref={textareaRef}
                    value={prompt}
                    rows={1}
                    placeholder="Something's on your mind?"
                    onChange={(e) => setPrompt(e.target.value)}
                  />
                </div>
                <div className='col'>
                  <button className='user-prompt-button' onClick={handlePrompt}>
                    {"->"}
                  </button>
                </div>
              </div>

              <div className='text-center'>
                <NewConversationButton />
              </div>
            </div>
          </div>
        }

        { loading &&
          <div className='row'>
            <div className='col assistant-response-bubble'>
              <div className="spinner-border" role="status">
                <span className="visually-hidden">Loading...</span>
              </div>
            </div>
          </div>
        }
      </div>
    </div>
  );
}

function UserMessage({ message }) {
  return (
    <div className='row'>
      <div className='col'/>
      <div className='col-10 user-message-bubble'>
        <p>
          {message.text}
        </p>
      </div>
    </div>
  );
}

function AssistantMessage({ index, message }) {
  const answerInTextId = `answerInText_${index}`;
  const [isCollapsed, setIsCollapsed] = useState(true);

  return (
    <div className='row'>
      <div className='col-10 assistant-response-bubble'>
        <audio src={message.speechClipUrl} controls />

        <p className='collapse' id={answerInTextId}>
          {message.text}
        </p>

        <button className='show-text-button'
                type='button'
                data-bs-toggle='collapse'
                data-bs-target={`#${answerInTextId}`}
                aria-expanded='false'
                onClick={_ => setIsCollapsed(!isCollapsed)}>
          {isCollapsed ? 'Show Text' : 'Hide'}
        </button>
      </div>
    </div>
  );
}

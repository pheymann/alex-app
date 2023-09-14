import {useState, useEffect, useRef} from 'react';
import { Link } from 'react-router-dom';
import './Home.css';
import NewConversationButton from './NewConversationButton';
import BasicPage from './BasicPage';
import { logError, pushLogMessage } from './logger';

export default function Home({ awsContext }) {
  const [conversations, setConversations] = useState([]);

  const logEntriesRef = useRef([]);

  useEffect(() => {
    const token = awsContext.token;

    fetch(`/api/conversation/list`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    })
      .then(response => response.text())
      .then(rawData => {
        pushLogMessage(logEntriesRef, { level: 'debug', message: rawData });

        const json = JSON.parse(rawData);
        setConversations(json);
      })
      .catch(error => {
        logError({ token, error, logEntriesRef: logEntriesRef});
        alert('Error getting conversations:\n' + error);
      });
  }, [awsContext.token, awsContext.userUUID]);

  return (
    <BasicPage awsContext={awsContext}>
      <div className='container container-limited-width'>
        <div className='row'>
          <div className='col text-center'>
            <NewConversationButton className='home-new-conversation-button' />
          </div>
        </div>

          { conversations &&
            conversations.map((conversation, index) => {
              const key = `${conversation.id}_${index}`;

              return (
                <div key={key} className='row'>
                  <div className='col'>
                    <Link className='conversation-link' to={`/conversation/${conversation.id}`}>
                      {conversation.metadata.artContext}
                    </Link>
                  </div>
                </div>
              );
            })
          }
      </div>
    </BasicPage>
  );
}

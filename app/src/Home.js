import {useState, useEffect, useRef} from 'react';
import { Link } from 'react-router-dom';
import './Home.css';
import NewConversationButton from './NewConversationButton';
import BasicPage from './BasicPage';
import { logError, pushLogMessage } from './logger';

export default function Home({ awsFetch, signOut }) {
  const [conversations, setConversations] = useState([]);

  const logEntriesRef = useRef([]);

  useEffect(() => {
    awsFetch.call(`/api/conversation/list`, {
      method: 'GET',
    })
      .then(rawData => {
        pushLogMessage(logEntriesRef, { level: 'debug', message: rawData });

        const json = JSON.parse(rawData);
        setConversations(json);
      })
      .catch(error => {
        logError({ awsFetch, error, logEntriesRef: logEntriesRef});
        // TODO - show error to user
        // alert('Error getting conversations:\n' + error);
      });
  }, [awsFetch]);

  return (
    <BasicPage awsFetch={ awsFetch } signOut={ signOut } >
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

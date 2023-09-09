import {useState, useEffect} from 'react';
import { Link } from 'react-router-dom';
import './Home.css';
import NewConversationButton from './NewConversationButton';
import BasicPage from './BasicPage';

export default function Home({ awsContext }) {
  const [conversations, setConversations] = useState([]);

  useEffect(() => {
    fetch(`/api/conversation/list`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${awsContext.token}`,
      },
    })
      .then(response => response.json())
      .then(data => {
        setConversations(data);
      })
      .catch(error => {
        console.log(error);
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

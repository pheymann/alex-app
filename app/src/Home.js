import { Auth } from 'aws-amplify';
import {useState, useEffect} from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './Home.css';
import NewConversationButton from './NewConversationButton';

export default function Home({ awsContext }) {
  const [conversations, setConversations] = useState([]);
  const navigate = useNavigate();

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
      });
  }, [awsContext.token, awsContext.userUUID]);

  return (
    <div className='container'>
      <h1>Let's talk about Art</h1>

      <button onClick={() => {
        Auth.signOut()
          .then(_ => navigate('/login'))
          .catch(err => console.log(err));
      }}>
        Sign Out
      </button>

      <div className='row'>
        <div className='col text-center'>
          <NewConversationButton />
        </div>
      </div>

        { conversations &&
          conversations.map((conversation, index) => {
            const key = `${conversation.id}_${index}`;

            return (
              <div className='row'>
                <div className='col'>
                  <Link key={key} className='conversation-link' to={`/conversation/${conversation.id}`}>
                    {conversation.metadata.artPiece} by {conversation.metadata.artistName}
                  </Link>
                </div>
              </div>
            );
          })
        }
    </div>
  );
}

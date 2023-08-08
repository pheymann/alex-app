import {useState, useEffect} from 'react';
import { Link } from 'react-router-dom';

export default function Home({ awsContext }) {
  const [conversations, setConversations] = useState([]);

  useEffect(() => {
    fetch(`/api/conversation/list`, {
      method: 'GET',
      headers: {
        'User-UUID': awsContext.userUUID,
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
    <div>
      <h1>Let's talk about Art</h1>

      <Link to={'/conversation/new'}>
        Start a new conversation
      </Link>

        { conversations &&
          conversations.map((conversation, index) => {
            const key = `${conversation.id}_${index}`;

            return (
              <div key={key}>
                <h2>
                  <Link to={`/conversation/${conversation.id}`}>
                    {conversation.metadata.artPiece}
                  </Link>
                </h2>
                <p>by {conversation.metadata.artistName}</p>
              </div>
            );
          })
        }
    </div>
  );
}

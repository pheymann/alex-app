import { Routes, Route, useNavigate } from 'react-router-dom';
import { Amplify, Auth } from 'aws-amplify';
import awsExports from './aws-exports';
import Home from './Home';
import Login from './Login';
import Conversation from './converation/Conversation';
import { useEffect, useState } from 'react';
import './App.css';

Amplify.configure({
  Auth: {
      region: awsExports.REGION,
      userPoolId: awsExports.USER_POOL_ID,
      userPoolWebClientId: awsExports.USER_POOL_APP_CLIENT_ID
  }
})

export default function App() {
  const [loading, setLoading] = useState(true);
  const [awsContext, setAwsContext] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const awsSession = await Auth.currentSession();

        setAwsContext({
          awsSession,
          token: awsSession.getIdToken().getJwtToken(),
          signOut: () => Auth.signOut(),
        });
        setLoading(false);
      } catch (error) {
        setLoading(false);
        navigate('/login');
      }
    };

    checkAuth();
  }, [navigate]);

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
    <Routes>
      <Route path='/login' element={<Login />} />
      <Route exact path='/' element={<Home awsContext={ awsContext } />} />
      <Route path='/conversation/:id' element={<Conversation awsContext={ awsContext } />} />
    </Routes>
  );
}

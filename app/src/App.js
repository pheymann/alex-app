import { Routes, Route, useNavigate } from 'react-router-dom';
import { Amplify, Auth } from 'aws-amplify';
import awsExports from './aws-exports';
import Home from './Home';
import ArtConversation from './ArtConversation';
import Login from './Login';
import { useEffect, useState } from 'react';

Amplify.configure({
  Auth: {
      region: awsExports.REGION,
      userPoolId: awsExports.USER_POOL_ID,
      userPoolWebClientId: awsExports.USER_POOL_APP_CLIENT_ID
  }
})

export default function App() {
  const [awsContext, setAwsContext] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const awsSession = await Auth.currentSession();
        const user = await Auth.currentAuthenticatedUser();

        setAwsContext({
          awsSession,
          user,
          userUUID: user.attributes.sub,
          token: awsSession.getIdToken().getJwtToken(),
        });
      } catch (error) {
        console.log(error);
        navigate('/login');
      }
    };

    checkAuth();
  }, [navigate]);

  if (!awsContext) {
    return(
      <div className="spinner-border" role="status">
        <span className="visually-hidden">Loading...</span>
      </div>
    );
  }

  return (
    <Routes>
      <Route path='/login' element={<Login />} />
      <Route exact path='/' element={<Home awsContext={ awsContext } />} />
      <Route path='/conversation/:id' element={<ArtConversation awsContext={ awsContext } />} />
    </Routes>
  );
}

import { Routes, Route, useNavigate } from 'react-router-dom';
import { Amplify, Auth } from 'aws-amplify';
import awsExports from './aws-exports';
import Home from './Home';
import ArtConversation from './ArtConversation';
import Login from './Login';
import { useEffect, useState } from 'react';

const isProduction = process.env.NODE_ENV === 'production';

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

    if (isProduction) {
      checkAuth();
    } else {
      setAwsContext({
        awsSession: null,
        user: {
          username: 'test',
          attributes: {
            sub: '1',
          },
        },
        userUUID: '1',
        token: 'test',
      });
    }
  }, [navigate]);

  return (
    <Routes>
      <Route path='/login' element={<Login />} />
      <Route exact path='/' element={<Home awsContext={ awsContext } />} />
      <Route path='/conversation/:id' element={<ArtConversation awsContext={ awsContext } />} />
    </Routes>
  );
}

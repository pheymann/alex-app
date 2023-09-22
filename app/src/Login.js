import { Authenticator } from '@aws-amplify/ui-react';
import { Link } from 'react-router-dom';

export default function Login() {
  return(
    <Authenticator>
      {_ =>{
        return(
          <Link to='/'>
            Back to Home
          </Link>
        );
      }}
    </Authenticator>
  );
}

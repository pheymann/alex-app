import { errorAlertMessage } from '../ErrorAlert';
import './ErrorField.css';

export default function ErrorField({ errorCode }) {
  const errorMessage = errorAlertMessage(errorCode);

  return (
    <div className='row'>
      <div className='col-10 error-field alert alert-warning'>
        { errorMessage }
      </div>
      <div className='col'>
      </div>
    </div>
  );
}

import { errorAlertMessage } from '../ErrorAlert';
import './ErrorField.css';

export default function ErrorField({ errorCode, i18n }) {
  const errorMessage = errorAlertMessage(errorCode, i18n);

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

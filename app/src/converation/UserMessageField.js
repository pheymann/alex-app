import './UserMessageField.css';

export default function UserMessageField({ message }) {
  return (
    <div className='row'>
      <div className='col'/>
      <div className='col-10 user-message-field'>
        <p>
          {message.text}
        </p>
      </div>
    </div>
  );
}

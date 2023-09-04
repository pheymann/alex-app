import { useState } from "react";
import './AssistantResponseField.css';

export default function AssistantResponseField({ index, message }) {
  const answerInTextId = `answerInText_${index}`;
  const [isCollapsed, setIsCollapsed] = useState(true);

  if (message.speechClipIsExpired) {
    return (
      <div className='row'>
        <div className='col-10 assistant-response-field'>
          <p>
            {message.text}
          </p>
        </div>
        <div className='col'>
        </div>
      </div>
    );
  }

  return (
    <div className='row'>
      <div className='col-10 assistant-response-field'>
        <audio src={message.speechClipUrl} controls />

        <p className='collapse' id={answerInTextId}>
          {message.text}
        </p>

        <button className='show-text-button'
                type='button'
                data-bs-toggle='collapse'
                data-bs-target={`#${answerInTextId}`}
                aria-expanded='false'
                onClick={_ => setIsCollapsed(!isCollapsed)}>
          {isCollapsed ? 'Show Text' : 'Hide'}
        </button>
      </div>
      <div className='col'>
      </div>
    </div>
  );
}

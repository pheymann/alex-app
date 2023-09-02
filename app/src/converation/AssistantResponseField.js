import { useState } from "react";
import './AssistantResponseField.css';

export function LoadingAssistantResponseField() {
  return (
    <div className='row'>
      <div className='col assistant-response-field'>
        <div className="spinner-border" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    </div>
  );
}

export function AssistantResponseField({ index, message }) {
  const answerInTextId = `answerInText_${index}`;
  const [isCollapsed, setIsCollapsed] = useState(true);

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
    </div>
  );
}

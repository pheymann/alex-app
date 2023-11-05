import { useState } from "react";
import './AssistantResponseField.css';

export default function AssistantResponseField({ index, message, i18n }) {
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
          {isCollapsed ? i18n.conversation.assistantResponse.show : i18n.conversation.assistantResponse.hide }
        </button>
      </div>
      <div className='col'>
      </div>
    </div>
  );
}

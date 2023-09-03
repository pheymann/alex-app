import { useEffect, useRef } from "react";
import "./PromptField.css";

export function PromptField(props) {
  const textareaRef = useRef(null);

  const resizeTextArea = () => {
    textareaRef.current.style.height = "auto";
    textareaRef.current.style.height = textareaRef.current.scrollHeight + "px";
  };

  useEffect(resizeTextArea, [props.value]);

  return (
    <div className='prompt-field'>
      { props.children }

      <div className='row'>
        <div className='col col-10'>
          <textarea
            className='prompt-field-textarea'
            ref={ textareaRef }
            rows={ 1 }
            value={ props.value }
            placeholder={ props.placeholder }
            onChange={(e) => props.onChangeValue(e.target.value)}
          />
        </div>
        <div className='col'>
          <button className='prompt-field-button' onClick={ props.onSubmit }>
            { "->" }
          </button>
        </div>
      </div>
    </div>
  );
}

export function LoadingPromptField() {
  return (
    <div className='row'>
      <div className='col d-flex prompt-field justify-content-center'>
        <div className="spinner-border" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    </div>
  );
}

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
    <div className='d-flex justify-content-center'>
      <div className='prompt-field'>
        <div className='container'>
          { props.children }

          <div className='row'>
            <div className='col col-9'>
              <textarea
                className='prompt-field-textarea'
                ref={ textareaRef }
                rows={ 1 }
                value={ props.value }
                placeholder={ props.placeholder }
                maxLength={ props.maxLength ? props.maxLength : 500 }
                autoFocus={ true }
                onChange={(e) => props.onChangeValue(e.target.value)}
              />
              <p className="prompt-field-remaining-characters text-end">{ props.value.length } / { props.maxLength }</p>
            </div>
            <div className='col-3'>
              <button className='prompt-field-button d-flex justify-content-center' onClick={ props.onSubmit }>
                Send
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export function LoadingPromptField() {
  return (
    <div className='d-flex justify-content-center'>
      <div className="prompt-field d-flex justify-content-center">
        <div className="spinner-border" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    </div>
  );
}

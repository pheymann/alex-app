import React from 'react';
import ReactDOM from 'react-dom/client';

import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.bundle.min";

import './index.css';
import { App } from './App';
import { BrowserRouter } from 'react-router-dom';
import { AwsFetch } from './awsfetch';
import { Language } from './language';
import { Auth } from 'aws-amplify';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <BrowserRouter>
      <App  validateSession={ () => Auth.currentSession() }
            buildAwsFetch={ (language) => new AwsFetch(language) }
            defaultLanguage={ Language.German }
       />
    </BrowserRouter>
  </React.StrictMode>
);

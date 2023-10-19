import React from 'react';
import ReactDOM from 'react-dom/client';

import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.bundle.min";

import './index.css';
import { App, defaultLoadAwsCtx } from './App';
import { BrowserRouter } from 'react-router-dom';
import { AwsFetch } from './awsfetch';
import { Language } from './language';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <BrowserRouter>
      <App  loadAwsCtx={ () => defaultLoadAwsCtx() }
            buildAwsFetch={ (awsContext, language) => new AwsFetch(awsContext, language) }
            defaultLanguage={ Language.German }
       />
    </BrowserRouter>
  </React.StrictMode>
);

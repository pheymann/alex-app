import React from 'react';
import ReactDOM from 'react-dom/client';

import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.bundle.min";

import './index.css';
import { App, defaultLoadAwsCtx } from './App';
import { BrowserRouter } from 'react-router-dom';
import { AwsFetch } from './awsfetch';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <BrowserRouter>
      <App  loadAwsCtx={ () => defaultLoadAwsCtx() }
            buildAwsFetch={ (awsContext) => new AwsFetch(awsContext) } />
    </BrowserRouter>
  </React.StrictMode>
);

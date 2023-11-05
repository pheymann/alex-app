import { render, waitFor } from '@testing-library/react';
import { App } from '../../App';
import { MemoryRouter } from 'react-router-dom';
import yaml from 'js-yaml';
import fs from 'fs';
import { AwsFetch } from '../../awsfetch';
import { fail } from 'assert';
import { Language } from '../../language';

export function runContract(contractPath, assertFn, userInteractionFn = async () => {}) {
  const contract = loadContract(contractPath);
  const mock = new AwsFetch(Language.English, new MockSessionTokenLoader(contract.authorizationToken), new MockFetchFn(contract.callChain))

  test("Case: " + contract.name, async () => {
    await render(
      <MemoryRouter initialEntries={[contract.view]}>
        <App  validateSession={ () => mockLoadAwsCtx(true) }
              buildAwsFetch={ (_) => mock }
              defaultLanguage={ Language.English }
        />,
      </MemoryRouter>
    );

    await userInteractionFn();

    await waitFor(() => assertFn(contract.app, mock.fetchFn.requestAndResponseCallCounter));
  });
}

function loadContract(contractPath) {
  const contract = yaml.load(
    fs.readFileSync(`./../cdc/${contractPath}`, 'utf8')
  );

  // adjust header structure to align with fetch api
  contract.callChain.forEach(requestAndResponse => {
    const headerObj = {};
    requestAndResponse.request.headers.forEach(pair => {
      headerObj[pair.name] = pair.value;
    });
    requestAndResponse.request.headerObj = headerObj;
  });

  return contract;
}

class MockFetchFn {

  constructor(requestsAndResponses) {
    this.requestsAndResponsesMap = new Map();
    this.requestAndResponseCallCounter = new Map();

    requestsAndResponses.forEach(requestAndResponse => {
      this.requestsAndResponsesMap.set(requestAndResponse.request.uri, requestAndResponse);
      this.requestAndResponseCallCounter.set(requestAndResponse.request.uri, 0);
    });
  }

  apply(uri, { method, headers, body }) {
    const requestAndResponse = this.requestsAndResponsesMap.get(uri);
    if (requestAndResponse === undefined) {
      fail(`No request and response found for uri: ${uri}\n${body}`);
    }

    expect(method).toEqual(requestAndResponse.request.method);
    expect(headers).toEqual(requestAndResponse.request.headerObj);

    const callCounter = this.requestAndResponseCallCounter.get(uri);
    this.requestAndResponseCallCounter.set(uri, callCounter + 1);

    return Promise.resolve({
      status: requestAndResponse.response.status,
      text: () => JSON.stringify(requestAndResponse.response.body),
      json: () => requestAndResponse.response.body,
      ok: requestAndResponse.response.status < 300,
    });
  }

  poll(uri, handleSuccess, handleError, props) {
    return this.apply(uri, props).then(response => {
      if (response.status < 300) {
        handleSuccess(response, () => {});
      }
      else {
        handleError(response);
      }
    })
  }
}

class MockSessionTokenLoader {

  constructor(token) {
    this.token = token;
  }

  load() {
    return Promise.resolve(this.token);
  }
}

async function mockLoadAwsCtx(isLoggedIn) {
  if (isLoggedIn) {
      return () => Promise.resolve({
      token: "ignore",
      signOut: () => {},
    });
  }

  return () => Promise.reject({
    message: "not logged in",
  });
}

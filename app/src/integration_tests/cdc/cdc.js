import { render, waitFor } from '@testing-library/react';
import { App } from '../../App';
import { MemoryRouter } from 'react-router-dom';
import yaml from 'js-yaml';
import fs from 'fs';
import { AwsFetch } from '../../awsfetch';
import { fail } from 'assert';

export function runContract(contractPath, assertFn, userInteractionFn = () => {}) {
  const contract = loadContract(contractPath);
  const mock = new AwsFetch({ token: contract.authorizationToken }, new MockFetchFn(contract.callChain))

  test("Case: " + contract.name, async () => {
    render(
      <MemoryRouter initialEntries={[contract.view]}>
        <App loadAwsCtx={ () => mockLoadAwsCtx(true) } buildAwsFetch={ (_) => mock } />,
      </MemoryRouter>,
    );

    await waitFor(() => {
      userInteractionFn();
    });

    await waitFor(() => {
      assertFn(contract.app, mock.fetchFn.requestAndResponseCallCounter);
    })
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

  apply(uri, { method, headers, _ }) {
    const requestAndResponse = this.requestsAndResponsesMap.get(uri);
    if (requestAndResponse === undefined) {
      fail(`No request and response found for uri: ${uri}`);
    }

    expect(method).toEqual(requestAndResponse.request.method);
    expect(headers).toEqual(requestAndResponse.request.headerObj);

    const callCounter = this.requestAndResponseCallCounter.get(uri);
    this.requestAndResponseCallCounter.set(uri, callCounter + 1);

    return Promise.resolve({
      status: requestAndResponse.response.status,
      text: () => JSON.stringify(requestAndResponse.response.body),
    });
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

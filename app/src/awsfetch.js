
export class AwsFetch {

  constructor(awsContext, fetchFn = new FetchWrapper()) {
    this.token = awsContext.token;
    this.fetchFn = fetchFn;
  }

  call(uri, { method, headers, body }) {
    if (headers === undefined) {
      headers = {};
    }
    headers['Authorization'] = `Bearer ${this.token}`;

    return this.fetchFn.apply(uri, {
      method: method,
      headers: headers,
      body: body,
    }).then(response => {
      if (response.status < 300) {
        return response.text();
      }
      else {
        throw new Error(`Request failed with status ${response.status}`);
      }
    });
  };
}

class FetchWrapper {
  apply(uri, props) {
    return fetch(uri, props);
  }
}

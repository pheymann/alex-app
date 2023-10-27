import { encodeLanguage } from "./language";

export class AwsFetch {

  constructor(awsContext, language, fetchFn = new FetchWrapper()) {
    this.token = awsContext.token;
    this.language = language;
    this.fetchFn = fetchFn;
  }

  call(uri, { method, headers, body }) {
    if (headers === undefined) {
      headers = {};
    }
    headers['Authorization'] = `Bearer ${this.token}`;
    headers['Accept-Language'] = encodeLanguage(this.language);

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

  callResponse(uri, { method, headers, body }) {
    if (headers === undefined) {
      headers = {};
    }
    headers['Authorization'] = `Bearer ${this.token}`;
    headers['Accept-Language'] = encodeLanguage(this.language);

    return this.fetchFn.apply(uri, {
      method: method,
      headers: headers,
      body: body,
    }).then(response => {
      if (response.status < 300) {
        return response;
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

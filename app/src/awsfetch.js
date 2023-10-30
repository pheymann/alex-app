import { Auth } from "aws-amplify";
import { encodeLanguage } from "./language";

export class AwsFetch {

  constructor(language, tokenLoader = new SessionTokenLoader(), fetchFn = new FetchWrapper()) {
    this.language = language;
    this.tokenLoader = tokenLoader;
    this.fetchFn = fetchFn;
  }

  async call(uri, { method, headers, body }) {
    if (headers === undefined) {
      headers = {};
    }

    const token = await this.tokenLoader.load();
    headers['Authorization'] = `Bearer ${token}`;
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

  async callResponse(uri, { method, headers, body }) {
    if (headers === undefined) {
      headers = {};
    }

    const token = await this.tokenLoader.load();
    headers['Authorization'] = `Bearer ${token}`;
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

class SessionTokenLoader {
  async load() {
    const awsSession = await Auth.currentSession();
    return awsSession.getIdToken().getJwtToken();
  }
}

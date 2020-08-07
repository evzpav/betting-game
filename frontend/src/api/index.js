import axios from "axios";
import config from "../config";

function getTest() {
  return baseRequest("GET", config.apiUrl, `/test`);
}

function postJoinGame(payload) {
    return baseRequest("POST", config.apiUrl, `/join`, payload);
  }

function baseRequest(method, url, uri, data = "") {
  return axios({
    method,
    timeout: 1000 * 5,
    url: `http://${url}/api${uri}`,
    data,
    headers: {
      "Content-Type": "application/json",
    },
  });
}

export { getTest, postJoinGame };
import axios from "axios";
import config from "../config";

function postJoinGame(payload) {
    return baseRequest("POST", config.apiUrl, `/game/join`, payload);
}

function getRankingSnapshot() {
  return baseRequest("GET", config.apiUrl, `/ranking/snapshot`);
}

function getGameSnapshot() {
  return baseRequest("GET", config.apiUrl, `/game/snapshot`);
}

function newWebsocket(){
  const wsId = uuidv4();
  const protocol = config.protocol === "https:" ? "wss:":"ws:";

  return  new WebSocket(`${protocol}//${config.apiUrl}/api/ws?id=${wsId}`);
}

function baseRequest(method, url, uri, data = "") {
  return axios({
    method,
    timeout: 1000 * 5,
    url: `${config.protocol}//${url}/api${uri}`,
    data,
    headers: {
      "Content-Type": "application/json",
    },
  });
}


function uuidv4() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

export { postJoinGame, getRankingSnapshot,getGameSnapshot, newWebsocket };
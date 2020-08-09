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
  return  new WebSocket("ws://" + config.apiUrl + "/api/ws");
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

export { postJoinGame, getRankingSnapshot,getGameSnapshot, newWebsocket };
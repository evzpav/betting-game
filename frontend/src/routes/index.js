import Home from "../components/Home";
import Leaderboard from "../components/Leaderboard";
import Rules from "../components/Rules";

export default [
  {
    path: "/",
    component: Home,
  },
  {
    path: "/leaderboard",
    component: Leaderboard,
  },
  {
    path: "/rules",
    component: Rules,
  }
];

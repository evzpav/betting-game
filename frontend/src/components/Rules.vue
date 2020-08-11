<template>
  <div class="content">
    <h2>Starting the game</h2>
    <ul>
      <li>2 players are needed to start the game (1 browser tab one player)</li>
    </ul>
    <ol>
      <li>On the initial page click "Play"</li>
      <li>Then fill your name</li>
      <li>Click on 2 numbers of your choice (1 to 10)</li>
      <li>Click "Join"</li>
      <li>Open another tab/window of this link in your browser</li>
      <li>Repeat in this other tab the steps 1 to 4</li>
      <li>Game will start</li>
    </ol>
    <ul>
      <li>When new players join in a running game they will only observe until next game starts, then they will join automatically.</li>
    </ul>

    <h2>Game Mechanics:</h2>
    <ul>
      <li>New random number generated every 1 second</li>
      <li>Maximum 30 rounds</li>
      <li>Ends if any player reaches exactly 21 points</li>
      <li>10 seconds of interval between games</li>
    </ul>

    <h2>Score:</h2>
    <ul>
      <li>
        Exact match
        <ul>
          <li>+5 points if there is exact match with one of the chosen number</li>
        </ul>
      </li>
      <li>
        Inside bounds
        <ul>
          <li>+5-(upper bound - lower bound); e.g.: 3 and 8 chosen, 7 is the round number, then 5-(8-3)=0</li>
        </ul>
      </li>
      <li>
        Out of bounds
        <ul>
          <li>-1 point; e.g.: 3 and 8 chosen, 9 is the round number</li>
        </ul>
      </li>
    </ul>

    <h2>Winner criteria (in this order):</h2>
    <ul>
      <li>Points</li>
      <li>Highest upper number chosen</li>
      <li>Highest lower number chosen</li>
      <li>Name ascending</li>
    </ul>

    <router-link id="back-btn" to="/" tag="button" class="button is-danger">BACK</router-link>
  </div>
</template>

<style scoped>
#back-btn {
  margin-top: 20px;
}

.content {
  font-size: 13px;
  margin: auto;
  padding: 20px;
}
</style>
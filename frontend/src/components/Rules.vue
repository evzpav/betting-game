<template>
  <div class="content">
    <p>Players chooses 2 numbers from 1 to 10 (including 10).</p>

    <p>Game:</p>
    <ul>
      <li>New random number generated every 1 second</li>
      <li>Maximum 30 rounds</li>
      <li>Ends if any player reaches exactly 21 points</li>
      <li>10s interval between games</li>
    </ul>

    <p>Score:</p>
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

    <p>Winner criteria (in this order):</p>
    <ul>
      <li>Points</li>
      <li>Highest upper number chosen</li>
      <li>Highest lower number chosen</li>
      <li>Name ascending</li>
    </ul>
  </div>
</template>

<style scoped>
</style>
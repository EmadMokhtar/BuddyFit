<template>
  <div id="app">
    <div class="header-container">
      <h1>BuddyFit Bot Chat</h1>
      <h2> How can I help you today champion?</h2>
    </div>
    <div class="chat-wrapper">
      <Chat />
    </div>
  </div>
</template>

<script>
import Chat from './components/Chat.vue';

export default {
  name: 'App',
  components: {
    Chat
  }
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: white;
  margin-top: 60px;
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: #181818;
}

.header-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.chat-wrapper {
  display: flex;
  justify-content: center;
  width: 100%;
}
</style>

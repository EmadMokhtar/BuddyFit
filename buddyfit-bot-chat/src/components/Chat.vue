<template>
  <div class="chat-container">
    <div class="messages">
      <div v-for="(message, index) in messages" :key="index" :class="{'user-message': message.user, 'ai-message': !message.user}">
        <span v-if="message.user" class="emoji">💪</span>
        <span v-else class="emoji">🤖</span>
        <p v-if="message.user">{{ message.text }}</p>
        <p v-else v-html="message.html"></p>
      </div>
    </div>
    <div class="input-container">
      <input v-model="input" @keyup.enter="sendMessage" placeholder="How can I help you champion?..." />
      <button @click="sendMessage">Send</button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import { marked } from 'marked';

export default {
  data() {
    return {
      input: '',
      messages: []
    };
  },
  methods: {
    async sendMessage() {
      if (this.input.trim() === '') return;

      this.messages.push({ text: this.input, user: true });
      const userInput = this.input;
      this.input = '';

      try {
        const response = await axios.post('http://localhost:8000/ask', { prompt: userInput });
        const markdownResponse = response.data.response;
        const htmlResponse = marked(markdownResponse);
        this.messages.push({ text: markdownResponse, html: htmlResponse, user: false });
      } catch (error) {
        console.error('Error sending message:', error);
      }
    }
  }
};
</script>

<style>
.chat-container {
  width: 80vw;
  height: 80vh;
  margin: 0 auto; /* Center the chat container */
  padding: 10px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.messages {
  flex: 1;
  padding: 10px;
  overflow-y: auto;
}

.user-message {
  text-align: right;
  background-color: #172859;
  padding: 10px;
  border-radius: 10px;
  margin: 5px 0;
  color: white;
}

.ai-message {
  text-align: left;
  background-color: #24262b;
  padding: 10px;
  border-radius: 10px;
  margin: 5px 0;
  color: white;
}

.input-container {
  display: flex;
  padding: 10px;
  border-top: 1px solid #ccc;
}

input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: large;
}

button {
  padding: 10px 20px;
  margin-left: 10px;
  border: none;
  background-color: #172859;
  color: white;
  border-radius: 4px;
  cursor: pointer;
}

ul {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

li {
  list-style-type: none;
}

.emoji {
  margin-right: 10px;
  font-size: x-large;
}
</style>

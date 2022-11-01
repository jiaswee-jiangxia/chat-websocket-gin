<template>
  <div class="common-layout">
    <el-container>
      <el-aside width="200px">
        <strong>Room list</strong>
        <p v-for="room in rooms" :key="room">
          {{ room }}
        </p>
      </el-aside>
      <el-container>
        <el-header>Header</el-header>
        <el-main>
          <el-empty v-if="messages.length==0" description="No message"></el-empty>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script>
var WS;
export default {
  name: 'App',
  data(){
    return {
      rooms:[],
      username:"default",
      messages:[],
      WS:null,
    }
  },
  mounted(){
    WS = new WebSocket("ws://localhost:20192/ws")
    WS.onopen = function(event) {
      console.log(event)
      console.log("Connected")
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>

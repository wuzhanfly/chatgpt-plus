<template>
  <div class="chat-line chat-line-right">
    <div class="chat-item">
      <div class="content" v-html="content"></div>
      <div class="triangle"></div>
    </div>

    <div class="chat-icon">
      <img :src="icon" alt="User"/>
    </div>
  </div>
</template>

<script>
import {defineComponent} from "vue"

export default defineComponent({
  name: 'ChatPrompt',
  props: {
    content: {
      type: String,
      default: '',
    },
    icon: {
      type: String,
      default: 'images/user-icon.png',
    }
  },
  data() {
    return {}
  },
})
</script>

<style lang="stylus" scoped>
.chat-line-right {
  justify-content: flex-end;

  .chat-icon {
    margin-left 5px;

    img {
      border-radius 50%;
    }
  }

  .chat-item {
    position: relative;
    padding: 0 5px 0 0;
    overflow: hidden;

    .triangle {
      width: 0;
      height: 0;
      border-top: 10px solid transparent;
      border-bottom: 10px solid transparent;
      border-left: 10px solid #223A34;
      position: absolute;
      right: 0;
      top: 10px;
    }

    .content {
      word-break break-word;
      padding: 12px 15px;
      background-color: #223A34;
      color var(--content-color);
      font-size: var(--content-font-size);
      border-radius: 5px;

      p {
        line-height 1.5
      }

      p:last-child {
        margin-bottom: 0
      }

      p:first-child {
        margin-top 0
      }
    }
  }
}
</style>
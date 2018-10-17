<style lang="scss" scoped>
.items {
  display: flex;
  flex-wrap: wrap;
  flex-direction: column;
  justify-content: center;
  height: 475px;
}

.list-area {
  z-index: 10;
  position: fixed;
background: #FFF;
border-color: #000;
margin-left: 10px;
}

.list-item {
list-style-type: none;
margin-left: 10px;
margin-right: 10px;
}

.item-selected {
  background: #41B883;
  color: #FFF
}

.items .item {
  flex: 1;
  box-sizing: border-box;
  margin: 8px 20px 8px 0;
  color: #171e42;
}

.item input {
  border: none;
  outline: none;
  color: blue;
  border-bottom: solid 1px black;
}

p.error {
  color: red;
  font-weight: bold;
}

.error .item input {
  color: red;
}
</style>

<template>
  <!-- POP UP - welcome dialog -->
  <div class="pop-up-box fullscreen">
    <div class="pop-up-box__content">
      <div class="pop-up-box__header">
        <h1>Enter recovery phrase</h1>
      </div>
      <p>
        Enter your 24-word recovery phrase to access your existing account:
      </p>
      <form class="items" :class="{error: errMsg}">
        <template v-for="(word, index) in words">
          <div class="item">
            {{ 1 + parseInt(index) + '.' }}
            <input
              type="text"
              v-model='words[index]'
              @keyup.down="calMatchListIndex(1)"
              @keyup.up="calMatchListIndex(-1)"
              @keyup.enter="updateValue(index)"
              @keyup.esc="resetMatchList()"
              @input="onChange($event.target.value, index)"
              />
            <div v-if="focusIdx === index" class="list-area">
            <div v-if="focusIdx === index && matchList.length" class="list-area">
              <li class="list-item" v-for="(item, matchIdx) in matchList">
                <p
                  :class="{'item-selected': isSelected(matchIdx)}"
                  v-on:click="updateValue(index, matchIdx)"
                  >
                  {{item}}
                </p>
              </li>
            </div>
          </div>
        </template>
      </form>
      </p>
      <p v-if="errMsg" class="error">Can not recover the account. Error: {{ errMsg }}</p>
      <div class="button">
        <button type="button" class="btn-default-fill" @click="backToWelcome">BACK</button>
        <button type="button" class="btn-default-fill" @click="recoverAccount">NEXT</button>
      </div>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import wordList from "./wordListEng";

export default {
  props: {
    maxCount: {
      type: Number,
      default: 5
    },
    minChars: {
      type: Number,
      default: 1
    },
    presetList: Array
  },

  methods: {
    recoverAccount() {
      let phrase = this.words.join(" ");
      axios
        .post("/api/account/phrase", {
          phrase: phrase
        })
        .then(resp => {
          this.$emit("accountRecovered");
        })
        .catch(err => {
          this.errMsg = err.response.data.message;
        });
    },

    backToWelcome() {
      this.$emit("backToWelcome");
    },

    calMatchListIndex(nextElement) {
      // do nothing for empty list
      if (!this.matchList.length) {
        return;
      }

      if (this.matchListIndex === null) {
        this.matchListIndex = 0;
        return;
      }

      // select to end of list, go back to first one
      if (this.matchListIndex + nextElement >= this.matchList.length) {
        this.matchListIndex = 0;
        return;
      }

      // select to beginning of list, go to last one
      if (this.matchListIndex + nextElement < 0) {
        this.matchListIndex = this.matchList.length - 1;
        return;
      }

      this.matchListIndex = this.matchListIndex + nextElement;
    },

    isSelected(index) {
      return index === this.matchListIndex;
    },

    resetMatchList() {
      this.matchListIndex = null;
      this.matchList = [];
    },

    clearMatchList() {
      this.matchList = [];
    },

    onChange(word, index) {
      // do nothing for emtpy list
      if (!word.length) {
        this.resetMatchList();
        return;
      }
      // know which input is typing
      this.focusIdx = index;

      // update only if words fix criteria
      if (word.length && word.length >= this.minChars) {
        const len = word.length;
        if (this.matchListIndex === null) {
          if (this.words[index] === "") {
            this.clearMatchList();
            return;
          }
          let matchData = wordList.filter(
            v => v.indexOf(word, 0) > -1 && v.substr(0, len) === word
          );
          this.matchList = matchData.slice(0, this.maxCount);
        }
      } else {
        this.clearMatchList();
      }
    },

    updateValue(index, matchIdx) {
      // do nothing for empty list
      if (!this.matchList.length) {
        return;
      }
      // default set to first element, if user has ever select, choose by that
      let targetIdx = this.matchListIndex || 0;

      // if user clicks by mouse, use matchIdx
      if (matchIdx !== undefined) {
        targetIdx = matchIdx;
      }

      this.words[index] = this.matchList[targetIdx];
      this.resetMatchList();
    }
  },

  data() {
    let words = [];
    for (let i = 0; i < 24; i++) {
      words[i] = "";
    }
    return {
      words: words,
      errMsg: "",
      matchList: [],
      matchListIndex: null,
      focusIdx: null
    };
  }
};
</script>

<style lang="scss" scoped>
.items {
  display: flex;
  flex-wrap: wrap;
  flex-direction: column;
  justify-content: center;
  height: 475px;
}

.items .item {
  flex: 1;
  box-sizing: border-box;
  margin: 8px 20px 8px 0;
  color: #171e42;
  position: relative;
}

.item input {
  border: none;
  outline: none;
  color: blue;
  border-bottom: solid 1px black;
}

.list-area {
  z-index: 10;
  background: #fff;
  border: 1px solid #dcdcdc;
  margin-left: 10px;
  border-radius: 5px;
  position: absolute;
  width: calc(100% - 80px);
  left: 5px;
  top: 21px;
}

.list-item {
  list-style-type: none;
}

.item-word {
  margin-left: 10px;
  margin-right: 10px;
}

.item-selected {
  height: 100%;
  background: #dcdcdc;
  color: #000;
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
  <div class="pop-up-box fullscreen" @click="resetMatchList()">
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
              @input="onChange($event.target.value, index)"
              />
            <div v-if="focusIdx === index && matchList.length" class="list-area">
              <li class="list-item" v-for="(item, matchIdx) in matchList">
                <div
                  :class="{'item-selected': isSelected(matchIdx)}"
                  @click="updateValue(index, matchIdx)"
                  >
                  <span class="item-word">
                    {{item}}
                  </span>
                </div>
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
    }
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

    onChange(word, index) {
      // do nothing for emtpy list
      if (!word.length) {
        this.resetMatchList();
        return;
      }
      // know which input is typing
      this.focusIdx = index;

      // update only if words fit criteria
      if (word.length && word.length >= this.minChars) {
        const len = word.length;
        let matchData = wordList.filter(
          v => v.indexOf(word, 0) > -1 && v.substr(0, len) === word
        );

        // remove list when exactly one match
        if (matchData.length === 1 && matchData[0] === word) {
          this.resetMatchList();
        } else {
          this.matchList = matchData.slice(0, this.maxCount);
        }
      } else {
        this.resetMatchList();
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

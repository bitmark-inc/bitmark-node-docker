<template>
  <!-- POP UP - welcome dialog -->
  <div class="pop-up-box fullscreen" >
    <div class="pop-up-box__content">
      <div class="pop-up-box__header">
        <h1>Welcome to Bitmark node</h1>
      </div>
      <p>
        Would you like this node to access an existing Bitmark account or generate a new one?
      </p>
      <div class="button">
        <button type="button" class="btn-default-fill" @click="existingAccount">EXISTING ACCOUNT</button>
        <button type="button" class="btn-default-fill" @click="createNewAccount">NEW ACCOUNT</button>
      </div>
    </div>
  </div>
</template>

<script>
  export default {
    methods: {
      createNewAccount() {
        this.$emit('createAccount')
      },

      existingAccount() {
        this.$emit('existingAccount')
      }
    }
  }
</script>

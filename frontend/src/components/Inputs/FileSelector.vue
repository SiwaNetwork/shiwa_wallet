<template>
  <div class="input-group">
    <input type="text"
      class="form-control"
      :disabled="true"
      aria-describedby="basic-addon2"
      :value=value />
    <span class="input-group-append">
      <p-button v-if="action == 'SelectFile'"
        tabIndex="-1" 
        @click.native="importKey" 
        class="btn"
        type="default">
        <span style="display: block;">
          ЗАГРУЗИТЬ
        </span>
      </p-button>
      <p-button v-if="action == 'SelectSaveFile'"
        tabIndex="-1" 
        @click.native="SelectDirToStoreKey" 
        class="btn"
        type="default">
        <span style="display: block;">
          ЗАГРУЗИТЬ
        </span>
      </p-button>
    </span>
  </div>
</template>

<script>
export default {
  props: {
    value: String,
    placeholder: String,
    action: String
  },

  methods: {
    importKey: function() {
      window.backend.WalletApplication.ImportKey().then(
        result => {
          if (result) {
            this.$store.commit('wallet/setKeystorePath', result);
          }
        }
      );
    },
    SelectDirToStoreKey: function() {
      window.backend.WalletApplication.SelectDirToStoreKey().then(
        result => {
          this.$store.commit('wallet/setKeystorePath', result);
        }
      );
    }
  }
}

</script>

<style scoped></style>
import Vue from "vue";

Vue.mixin({
  methods: {
    logout: function() {
      window.backend.WalletApplication.LogOut().then(txFinishedState => {
        if (txFinishedState) {
          let darkMode = this.$store.state.wallet.darkMode
          this.$store.dispatch('transaction/reset').then(() => {
            this.$store.dispatch('wallet/reset').then(() => {
              this.$store.dispatch('app/reset').then(() => {
                this.$router.push({
                  name: 'login',
                  params: {
                    message: "Пожалуйста, введите свои учетные данные ниже, чтобы получить доступ к вашему кошельку SHIWA",
                    darkMode: darkMode
                  }
                });
                return;                
              })
            })            
          })
        }
      }), (this.random = "1");
      return;
    }
  }
});
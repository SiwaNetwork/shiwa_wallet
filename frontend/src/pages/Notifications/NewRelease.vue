<template>
  <span>There's a new version of Shiwa Wallet available. Click this notification to upgrade to v.{{ this.$store.state.app.newRelease }}.</span>
</template>

<script>
  export default {
    name: 'newRelease'
  }
</script>

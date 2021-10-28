<template>
  <v-app>
    <AppBar :menu="menu" showBar/>
    <nuxt />
  </v-app>
</template>

<script>
import AppBar from '@/components/AppBar'
export default {
  components: { AppBar },
  data: () => ({
    menu: [

    ]
  })
}
</script>

<style>
body {
  margin: 0 !important;
}
</style>

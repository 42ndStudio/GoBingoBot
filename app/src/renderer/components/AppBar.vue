<template>
  <v-app-bar app fixed color="white" alt="Aguapanelazo" light v-show="showBar">
    <v-toolbar-title class="logo ml-2">
      <nuxt-link to="/"> GoBingo </nuxt-link>
    </v-toolbar-title>

    <v-spacer></v-spacer>

    <v-toolbar-items
      v-show="menu && this.$vuetify.breakpoint.width > 1100"
      class="mr-2"
    >
      <v-btn
        v-for="item in menu"
        :key="item.icon"
        :href="item.href"
        :to="item.link"
        :icon="item.icon ? true : false"
        :text="item.icon ? false : true"
      >
        {{ item.text }}
      </v-btn>
    </v-toolbar-items>
    <img style="max-width: 100%" src="~assets/logo-42nd-Studio.png" />
  </v-app-bar>
</template>

<script>

export default {
  props: {
    menu: Array,
    showBar: Boolean,
  },
  data: () => ({
    
  }),
  methods: {
    
  },
  mounted() {
    console.log(this.$route);
  },
  components: {  },
};
</script>

<style>
.logo {
  overflow: unset;
}
img {
  height: 55px;
}
</style>
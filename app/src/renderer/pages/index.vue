<template>
  <v-container>
    <v-row align="center" justify="center">
      <v-col cols="12">
        <GamesList></GamesList>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { remote } from "electron";
import SystemInformation from "@/components/SystemInformation.vue";
import GamesList from "@/components/GamesList.vue";

export default {
  components: {
    GamesList,
    SystemInformation,
  },
  data() {
    return {
      externalContent: "",
    };
  },
  methods: {
    openURL(url) {
      remote.shell.openExternal(url);
    },
  },
};
</script>

<style>

</style>

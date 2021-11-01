<template>
  <v-chip :color="currentModeComp[1]" :class="currentModeComp[2]" small>{{
    currentModeComp[0]
  }}</v-chip>
</template>

<script>
export default {
  props: {
    currentMode: String,
  },
  computed: {
    currentModeComp() {
      if (!this.currentMode) {
        return "N/A", "red";
      }
      switch (this.currentMode) {
        case "a":
          return ["Todas las casillas", "green", "white--text"];

        case "o":
          return ["Hacer la O", "orange", "white--text"];

        case "n":
          return ["Hacer la N", "orange", "white--text"];

        case "/":
          return ["Diagonal de Abajo a Arriba (/)", "orange", "white--text"];

        case "\\":
          return ["Diagonal de  Arriba a Abajo (/)", "orange", "white--text"];
      }
      //   "a", "c", "o", "n", "/", "\\", "l1", "l2", "l3", "l4", "l5"}
      if (this.currentMode[0] == "l") {
        if (isNaN(this.currentMode[1])) {
          return ["Columna letra " + this.currentMode[1].toUpperCase(), "purple", "white--text"];
        } else {
          return ["Fila " + this.currentMode[1], "purple", "white--text"];
        }
      }
      return [cm, "red", "white--text"];
    },
  },
};
</script>

<style>
</style>
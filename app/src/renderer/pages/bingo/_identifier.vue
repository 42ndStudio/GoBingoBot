<template>
  <v-container>
    <v-row align="center" justify="center">
      <v-col cols="12">
        <v-card v-if="bingo && bingo.BingoID">
          <v-card-title> Bingo {{ bingo.name }}</v-card-title>
          <v-card-subtitle>{{ $route.params['identifier'] }}</v-card-subtitle>
          <v-card-text>
            <b>Tableros: {{ bingo.BoardsSold }}</b><br>
            <b>Modo de Juego: <GameModeChip :currentMode="bingo.CurrentMode" /></b><br>
            <v-row class="mt-4">
              <v-col cols="4">
                <v-btn color="green accent-2" @click="reqGenerateBoard" small>
                  Generar Tablero
                </v-btn>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import local_api from "~/api/local_api";
import GameModeChip from "~/components/GameModeChip";

export default {
  components: {
    GameModeChip
  },
  data: () => ({
    bingo: false,
    generated: false
  }),
  created() {
    console.log(this.$route);
  },
  mounted() {
    console.log(this.$route);
    this.loadGame()
  },
  methods: {
    loadGame() {
      console.log("loadGame", this.$route.params['identifier']);
      try {
        local_api.api_get("game", this.$route.params['identifier']).then((respuesta) => {
          console.log(respuesta);
          this.bingo = respuesta;
        });
      } catch (error) {
        console.error("failed getting games", error);
      }
    },
    reqGenerateBoard() {
      console.log("reqGenerateBoard")
      try {
        local_api.api_post("game/generate", {
          BingoID: this.$route.params['identifier'],
        }).then((respuesta) => {
          console.log(respuesta);
          this.generated = respuesta;
        });
      } catch (error) {
        console.error("failed getting games", error);
      }
    },
    submitNewGame() {
      console.log("saving game");
      this.dialog_add = false;
    },
  },
};
</script>

<style>
</style>